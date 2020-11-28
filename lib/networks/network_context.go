/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package networks

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kurtosis-tech/kurtosis-go/lib/kurtosis_service"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type NetworkContext struct {
	kurtosisService *kurtosis_service.KurtosisService

	// The dirpath ON THE SUITE CONTAINER where the suite execution volume is mounted
	suiteExecutionVolumeDirpath string

	services map[ServiceID]ServiceNode
}
/*
Creates a new NetworkContext object with the given parameters.

Args:
	kurtosisService: The Docker manager that will be used for manipulating the Docker engine during test network modification.
	servicesDirpath: The dirpath where directories for each new service will be created to store file IO
*/
func NewNetworkContext(
		kurtosisService *kurtosis_service.KurtosisService,
		servicesRelativeDirpath string) *NetworkContext {
	return &NetworkContext{
		kurtosisService:             kurtosisService,
		services:                    make(map[ServiceID]ServiceNode),
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
	configurationId: The ID of the service configuration to use for creating the service.
	serviceId: The service ID that will be used to identify this node in the network.
	dependencies: A "set" of service IDs that the node being created will depend on - i.e., whose information the node-to-create
		needs to start up. If the node-to-create doesn't depend on any other services, the dependencies map should be
		empty (not nil).

Return:
	An AvailabilityChecker for checking when the new service is available and ready for use.
*/
func (networkCtx *NetworkContext) GetServiceBuilder(serviceId ServiceID, core *services.ServiceBuilderCore) (*services.ServiceBuilder, error) {
	if _, exists := networkCtx.services[serviceId]; exists {
		return nil, stacktrace.NewError("Service ID %s already exists in the networkCtx", serviceId)
	}

	logrus.Trace("Creating directory within test volume for service...")
	serviceDirname := fmt.Sprintf("%v-%v", serviceId, uuid.New().String())
	serviceDirpath := filepath.Join(TODO, networkCtx.suiteExecutionVolumeDirpath, serviceDirname)
	err := os.Mkdir(serviceDirpath, os.ModeDir)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred creating the new service's directory in the volume at filepath '%v'", serviceDirpath)
	}
	logrus.Tracef("Successfully created directory for service: %v", serviceDirpath)

	// TODO this is goofy
	mountServiceDirpath := filepath.Join(core.GetTestVolumeMountpoint(), serviceDirname)

	builder := services.NewServiceBuilder(core, serviceDirpath, networkCtx.kurtosisService)
	return builder, nil
}

/*
Gets the node information for the service with the given service ID.
*/
func (networkCtx *NetworkContext) GetService(serviceId ServiceID) (services.Service, error) {
	node, found := networkCtx.services[serviceId]
	if !found {
		return nil, stacktrace.NewError("No service with ID %v exists in the networkCtx", serviceId)
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
