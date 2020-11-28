/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package networks

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kurtosis-tech/kurtosis-go/lib/kurtosis_service"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

const (
	ipPlaceholder = "KURTOSISSERVICEIP"
)

/*
The identifier used for services with the network.
*/
type ServiceID string

/*
A package object containing the details that the NetworkContext is tracking about a node.
*/
type NetworkNode struct {
	// The user-defined interface for interacting with the node.
	// NOTE: this will need to be casted to the appropriate interface becaus Go doesn't yet have generics!
	Service services.Service

	// The Docker container ID running a given service
	ContainerID string
}


type NetworkContext struct {
	kurtosisService kurtosis_service.KurtosisService

	// The dirpath ON THE SUITE CONTAINER where the suite execution volume is mounted
	suiteExecutionVolumeDirpath string

	// Filepath to the services directory, RELATIVE to the root of the suite execution volume root
	servicesRelativeDirpath string

	services map[ServiceID]NetworkNode
}

/*
Creates a new NetworkContext object with the given parameters.

Args:
	kurtosisService: The Docker manager that will be used for manipulating the Docker engine during test network modification.
	servicesRelativeDirpath: The dirpath where directories for each new service will be created to store file IO, which
		is RELATIVE to the root of the suite execution volume!
*/
func NewNetworkContext(
		kurtosisService kurtosis_service.KurtosisService,
		servicesRelativeDirpath string) *NetworkContext {
	return &NetworkContext{
		kurtosisService:             kurtosisService,
		services:                    make(map[ServiceID]NetworkNode),
		suiteExecutionVolumeDirpath: servicesRelativeDirpath,
	}
}

// Gets the number of nodes in the network
func (networkCtx *NetworkContext) GetSize() int {
	return len(networkCtx.services)
}

/*
Adds a service to the network with the given service ID, created using the given configuration ID.

Args:
	serviceId: The service ID that will be used to identify this node in the network.
	initializer: The Docker container initializer that contains the logic for starting the service

Return:
	The new service
	An availability checker which can be used to wait until the service is available, if desired
*/
func (networkCtx *NetworkContext) AddService(
		serviceId ServiceID,
		initializer services.DockerContainerInitializer) (services.Service, services.AvailabilityChecker, error) {
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
	//  container so that it can do a "pre-registration" - register an IP address before actually starting the container
	if err := initializer.InitializeMountedFiles(osFiles); err != nil {
		return nil, nil, stacktrace.Propagate(err, "An error occurred initializing the files before service start")
	}
	logrus.Tracef("Successfully initialized files needed for service")

	logrus.Tracef("Creating start command for service...")
	startCmdArgs, err := initializer.GetStartCommand(mountFilepaths, ipPlaceholder)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "Failed to create start command")
	}
	logrus.Tracef("Successfully created start command for service")

	logrus.Tracef("Calling to Kurtosis API to create service...")
	dockerImage := initializer.GetDockerImage()
	ipAddr, containerId, err := networkCtx.kurtosisService.AddService(
		dockerImage,
		initializer.GetUsedPorts(),
		ipPlaceholder,
		startCmdArgs,
		make(map[string]string),
		initializer.GetTestVolumeMountpoint())
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "Could not add service for Docker image %v", dockerImage)
	}
	logrus.Tracef("Kurtosis API returned IP for new service: %v", ipAddr)

	logrus.Tracef("Getting service from IP...")
	service := initializer.GetServiceFromIp(ipAddr)
	logrus.Tracef("Successfully got service from IP")

	networkCtx.services[serviceId] = NetworkNode{
		Service:     service,
		ContainerID: containerId,
	}

	availabilityChecker := services.NewDefaultAvailabilityChecker(service)

	return service, availabilityChecker, nil
}

/*
Gets the node information for the service with the given service ID.
*/
func (networkCtx *NetworkContext) GetService(serviceId ServiceID) (services.Service, error) {
	node, found := networkCtx.services[serviceId]
	if !found {
		return nil, stacktrace.NewError("No service with ID %v exists in the network", serviceId)
	}

	return node.Service, nil
}

/*
Stops the container with the given service ID, and removes it from the network.
*/
func (networkCtx *NetworkContext) RemoveService(serviceId ServiceID, containerStopTimeoutSeconds int) error {
	nodeInfo, found := networkCtx.services[serviceId]
	if !found {
		return stacktrace.NewError("No service with ID %v found", serviceId)
	}

	logrus.Debugf("Removing service ID %v...", serviceId)
	delete(networkCtx.services, serviceId)

	// Make a best-effort attempt to stop the container
	err := networkCtx.kurtosisService.RemoveService(nodeInfo.ContainerID, containerStopTimeoutSeconds)
	if err != nil {
		logrus.Errorf(
			"The following error occurred stopping service ID %v with container ID %v; proceeding to stop other containers:",
			serviceId,
			nodeInfo.ContainerID)
		fmt.Fprintln(logrus.StandardLogger().Out, err)
	}
	logrus.Debugf("Successfully removed service ID %v", serviceId)
	return nil
}

/*
Makes a best-effort attempt to remove all the containers in the network, waiting for the given timeout and returning
	an error if the timeout is reached.

Args:
	containerStopTimeoutSeconds: How long to wait, in seconds, for each container to stop before force-killing it
*/
func (networkCtx *NetworkContext) RemoveAll(containerStopTimeoutSeconds int) error {
	for serviceId, _ := range networkCtx.services {
		networkCtx.RemoveService(serviceId, containerStopTimeoutSeconds)
	}
	return nil
}
