/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package networks

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kurtosis-tech/kurtosis-go/lib_core/kurtosis_service"
	"github.com/kurtosis-tech/kurtosis-go/lib_core/kurtosis_service/method_types"
	"github.com/kurtosis-tech/kurtosis-go/lib/client/artifact_id_provider"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
)

const (
	ipPlaceholder = "KURTOSISSERVICEIP"

	// This will alwyas resolve to the default partition ID (regardless of whether such a partition exists in the network,
	//  or it was repartitioned away)
	defaultPartitionId PartitionID = ""
)

type NetworkContext struct {
	// Mutex, to make this thread-safe
	mutex *sync.Mutex

	filesArtifactIdToGlobalArtifactIds map[services.FilesArtifactID]artifact_id_provider.ArtifactID

	kurtosisService kurtosis_service.KurtosisService

	// The dirpath ON THE SUITE CONTAINER where the suite execution volume is mounted
	suiteExecutionVolumeDirpath string

	// Filepath to the services directory, RELATIVE to the root of the suite execution volume root
	servicesRelativeDirpath string

	// The user-defined interfaces for interacting with the node.
	// NOTE: these will need to be casted to the appropriate interface becaus Go doesn't yet have generics!
	services map[services.ServiceID]services.Service
}


/*
Creates a new NetworkContext object with the given parameters.

Args:
	kurtosisService: The Docker manager that will be used for manipulating the Docker engine during test network modification.
	suiteExecutionVolumeDirpath: The path ON THE TEST SUITE CONTAINER where the suite execution volume is mounted
	servicesRelativeDirpath: The dirpath where directories for each new service will be created to store file IO, which
		is RELATIVE to the root of the suite execution volume!
	filesArtifactIdToGlobalArtifactId: Lookup table mapping files artifact IDs to global artifact IDs, for use when
		instantiating new services
*/
func NewNetworkContext(
		kurtosisService kurtosis_service.KurtosisService,
		suiteExecutionVolumeDirpath string,
		servicesRelativeDirpath string,
		filesArtifactIdToGlobalArtifactId map[services.FilesArtifactID]artifact_id_provider.ArtifactID) *NetworkContext {
	return &NetworkContext{
		mutex: &sync.Mutex{},
		filesArtifactIdToGlobalArtifactIds: filesArtifactIdToGlobalArtifactId,
		kurtosisService: kurtosisService,
		suiteExecutionVolumeDirpath: suiteExecutionVolumeDirpath,
		servicesRelativeDirpath: servicesRelativeDirpath,
		services: map[services.ServiceID]services.Service{},
	}
}

// Gets the number of nodes in the network
func (networkCtx *NetworkContext) GetSize() int {
	networkCtx.mutex.Lock()
	defer networkCtx.mutex.Unlock()

	return len(networkCtx.services)
}

/*
Adds a service to the network in the default partition with the given service ID, created using the given configuration ID.

NOTE: If the network has been repartitioned and the default partition hasn't been preserved, you should use
	AddServiceToPartition instead.

Args:
	serviceId: The service ID that will be used to identify this node in the network.
	initializer: The Docker container initializer that contains the logic for starting the service

Return:
	service: The new service
*/
func (networkCtx *NetworkContext) AddService(
		serviceId services.ServiceID,
		initializer services.DockerContainerInitializer) (services.Service, services.AvailabilityChecker, error) {
	// Mutex locked directly inside (we can't lock the mutex here because Go mutexes aren't reentrant)
	service, availabilityChecker, err := networkCtx.AddServiceToPartition(
		serviceId,
		defaultPartitionId,
		initializer)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "An error occurred adding the service to the network in the default partition")
	}
	return service, availabilityChecker, nil
}

/*
Adds a service to the network with the given service ID, created using the given configuration ID.

NOTE: If the network hasn't been repartitioned yet, the PartitionID should be an empty string to add to the default
	partition.

Args:
	serviceId: The service ID that will be used to identify this node in the network.
	partitionId: The partition ID to add the service to
	initializer: The Docker container initializer that contains the logic for starting the service

Return:
	service.Service: The new service
*/
func (networkCtx *NetworkContext) AddServiceToPartition(
		serviceId services.ServiceID,
		partitionId PartitionID,
		initializer services.DockerContainerInitializer) (services.Service, services.AvailabilityChecker, error) {
	networkCtx.mutex.Lock()
	defer networkCtx.mutex.Unlock()

	if _, exists := networkCtx.services[serviceId]; exists {
		return nil, nil, stacktrace.NewError("Service ID %s already exists in the network", serviceId)
	}

	serviceDirname := fmt.Sprintf("%v-%v", serviceId, uuid.New().String())
	serviceRelativeDirpath := filepath.Join(networkCtx.servicesRelativeDirpath, serviceDirname)

	logrus.Trace("Creating directory within test volume for service...")
	testSuiteServiceDirpath := filepath.Join(networkCtx.suiteExecutionVolumeDirpath, serviceRelativeDirpath)
	err := os.Mkdir(testSuiteServiceDirpath, os.ModeDir)
	if err != nil {
		return nil, nil, stacktrace.Propagate(
			err,
			"An error occurred creating the new service's directory in the volume at filepath '%v' on the testsuite",
			testSuiteServiceDirpath)
	}
	logrus.Tracef("Successfully created directory for service: %v", testSuiteServiceDirpath)

	mountServiceDirpath := filepath.Join(initializer.GetTestVolumeMountpoint(), serviceRelativeDirpath)

	logrus.Trace("Initializing files needed for service...")
	requestedFiles := initializer.GetFilesToMount()
	osFiles := make(map[string]*os.File)
	mountFilepaths := make(map[string]string)
	for fileId, _ := range requestedFiles {
		filename := fmt.Sprintf("%v-%v", fileId, uuid.New().String())
		testSuiteFilepath := filepath.Join(testSuiteServiceDirpath, filename)
		fp, err := os.Create(testSuiteFilepath)
		if err != nil {
			return nil, nil, stacktrace.Propagate(
				err,
				"Could not create new file for requested file ID '%v'",
				fileId)
		}
		defer fp.Close()
		osFiles[fileId] = fp
		mountFilepaths[fileId] = filepath.Join(mountServiceDirpath, filename)
	}
	// NOTE: If we need the IP address when initializing mounted files, we'll need to rejigger the Kurtosis API
	//  container so that it can do a "pre-registration" - dole out an IP address before actually starting the container
	if err := initializer.InitializeMountedFiles(osFiles); err != nil {
		return nil, nil, stacktrace.Propagate(err, "An error occurred initializing the files before service start")
	}
	logrus.Tracef("Successfully initialized files needed for service")

	logrus.Tracef("Creating files artifact mount dirpaths map...")
	filesArtifactMountDirpaths := map[artifact_id_provider.ArtifactID]string{}
	for filesArtifactId, mountDirpath := range initializer.GetFilesArtifactMountpoints() {
		globalArtifactId, found := networkCtx.filesArtifactIdToGlobalArtifactIds[filesArtifactId]
		if !found {
			return nil, nil, stacktrace.Propagate(
				err,
				"Service requested files artifact with ID '%v' to be mounted, but no" +
					"artifact with that ID was declared as needed in the test configuration",
				filesArtifactId)
		}
		filesArtifactMountDirpaths[globalArtifactId] = mountDirpath
	}
	logrus.Tracef("Successfully created files artifact mount dirpaths map")

	logrus.Tracef("Creating start command for service...")
	startCmdArgs, err := initializer.GetStartCommand(mountFilepaths, ipPlaceholder)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "Failed to create start command")
	}
	logrus.Tracef("Successfully created start command for service")

	logrus.Tracef("Calling to Kurtosis API to create service...")
	dockerImage := initializer.GetDockerImage()
	ipAddr, err := networkCtx.kurtosisService.AddService(
		string(serviceId),
		string(partitionId),
		dockerImage,
		initializer.GetUsedPorts(),
		ipPlaceholder,
		startCmdArgs,
		make(map[string]string),
		initializer.GetTestVolumeMountpoint(),
		filesArtifactMountDirpaths)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "Could not add service for Docker image %v", dockerImage)
	}
	logrus.Tracef("Kurtosis API returned IP for new service: %v", ipAddr)

	logrus.Tracef("Getting service from IP...")
	service := initializer.GetService(serviceId, ipAddr)
	logrus.Tracef("Successfully got service from IP")

	networkCtx.services[serviceId] = service

	availabilityChecker := services.NewDefaultAvailabilityChecker(serviceId, service)

	return service, availabilityChecker, nil
}

/*
Gets the node information for the service with the given service ID.
*/
func (networkCtx *NetworkContext) GetService(serviceId services.ServiceID) (services.Service, error) {
	networkCtx.mutex.Lock()
	defer networkCtx.mutex.Unlock()

	service, found := networkCtx.services[serviceId]
	if !found {
		return nil, stacktrace.NewError("No service with ID '%v' exists in the network", serviceId)
	}

	return service, nil
}

/*
Stops the container with the given service ID, and removes it from the network.
*/
func (networkCtx *NetworkContext) RemoveService(serviceId services.ServiceID, containerStopTimeoutSeconds int) error {
	networkCtx.mutex.Lock()
	defer networkCtx.mutex.Unlock()

	_, found := networkCtx.services[serviceId]
	if !found {
		return stacktrace.NewError("No service with ID %v found", serviceId)
	}

	logrus.Debugf("Removing service '%v'...", serviceId)
	delete(networkCtx.services, serviceId)

	// Make a best-effort attempt to stop the container
	err := networkCtx.kurtosisService.RemoveService(string(serviceId), containerStopTimeoutSeconds)
	if err != nil {
		return stacktrace.Propagate(err,
			"An error occurred removing service with ID '%v'",
			serviceId)
	}
	logrus.Debugf("Successfully removed service ID %v", serviceId)
	return nil
}

/*
Constructs a new repartitioner builder in preparation for a repartition.

Args:
	isDefaultPartitionConnectionBlocked: If true, when the connection details between two partitions aren't specified
		during a repartition then traffic between them will be blocked by default
 */
func (networkCtx NetworkContext) GetRepartitionerBuilder(isDefaultPartitionConnectionBlocked bool) *RepartitionerBuilder {
	// This function doesn't need a mutex lock because (as of 2020-12-28) it doesn't touch internal state whatsoever
	return newRepartitionerBuilder(isDefaultPartitionConnectionBlocked)
}

/*
Repartitions the network using the given repartitioner. A repartitioner builder can be constructed using the
	NewRepartitionerBuilder method of this network context object.
 */
func (networkCtx *NetworkContext) RepartitionNetwork(repartitioner *Repartitioner) error {
	networkCtx.mutex.Lock()
	defer networkCtx.mutex.Unlock()

	partitionServices := map[string]map[string]bool{}
	for partitionId, serviceIdSet := range repartitioner.partitionServices {
		serviceIdStrPseudoSet := map[string]bool{}
		for _, serviceId := range serviceIdSet.getElems() {
			serviceIdStr := string(serviceId)
			serviceIdStrPseudoSet[serviceIdStr] = true
		}
		partitionIdStr := string(partitionId)
		partitionServices[partitionIdStr] = serviceIdStrPseudoSet
	}

	serializablePartConns := map[string]map[string]method_types.SerializablePartitionConnection{}
	for partitionAId, partitionAConns := range repartitioner.partitionConnections {
		serializablePartAConns := map[string]method_types.SerializablePartitionConnection{}
		for partitionBId, unserializableConn := range partitionAConns {
			partitionBIdStr := string(partitionBId)
			serializableConn := makePartConnSerializable(unserializableConn)
			serializablePartAConns[partitionBIdStr] = serializableConn
		}
		partitionAIdStr := string(partitionAId)
		serializablePartConns[partitionAIdStr] = serializablePartAConns
	}

	serializableDefaultConn := makePartConnSerializable(repartitioner.defaultConnection)

	if err := networkCtx.kurtosisService.Repartition(partitionServices, serializablePartConns, serializableDefaultConn); err != nil {
		return stacktrace.Propagate(err, "An error occurred repartitioning the test network")
	}
	return nil
}

// ============================================================================================
//                                    Private helper methods
// ============================================================================================
func makePartConnSerializable(connection PartitionConnection) method_types.SerializablePartitionConnection {
	return method_types.SerializablePartitionConnection{
		IsBlocked: connection.IsBlocked,
	}
}
