package networks

import (
	"context"
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/lib/kurtosis_service"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
)

// TODO Rename this to ServiceTag
/*
The identifier used for services with the network.
 */
type ServiceID string

/*
A package object containing the details that the ServiceNetwork is tracking about a node.
 */
type ServiceNode struct {
	// The user-defined interface for interacting with the node.
	// NOTE: this will need to be casted to the appropriate interface becaus Go doesn't yet have generics!
	Service services.Service

	// The Docker container ID running a given service
	ContainerID string
}

/*
A package object containing the details of a particular service configuration, to give Kurtosis the implementation-specific
	details about how to interact with user-defined services.
 */
type serviceConfig struct {
	// The Docker image that will be used to launch nodes
	dockerImage string

	// The implementation that will be used for launching a Docker image of a node using this configuration
	initializerCore services.ServiceInitializerCore

	// The implementation that will be used for determining whether a node launched using this configuration is available
	availabilityCheckerCore services.ServiceAvailabilityCheckerCore
}


/*
A struct representing a network of services that will be used for a single test (commonly called the "test network"). This
	struct is the low-level access point for modifying the test network.
 */
type ServiceNetwork struct {
	// The Kurtosis service used for interacting with the Docker engine during test network manipulation
	kurtosisService *kurtosis_service.KurtosisService

	// A mapping of human-readable Service ID -> information about a node
	serviceNodes map[ServiceID]ServiceNode

	// A mapping of configuration ID -> configuration details
	configurations map[ConfigurationID]serviceConfig

	// The dirpath where the test volume is mounted on *the test suite container* (which is where this code will be running)
	testVolumeDirpath string
}

/*
Creates a new ServiceNetwork object with the given parameters.

Args:
	freeIpTracker: The IP tracker that will be used to provide IPs for new nodes added to the network.
	kurtosisService: The Docker manager that will be used for manipulating the Docker engine during test network modification.
	dockerNetworkName: The name of the Docker network this test network is running on.
	configurations: The configurations that are available for spinning up new nodes in the network.
	testVolume: The name of the Docker volume that will be mounted on all the nodes in the network.
	testVolumeDirpath: The dirpath that the test Docker volume is mounted on in the controller image (which will
		be running all the code here).
 */
func NewServiceNetwork(
			kurtosisService *kurtosis_service.KurtosisService,
			configurations map[ConfigurationID]serviceConfig,
			testVolumeDirpath string) *ServiceNetwork {
	return &ServiceNetwork{
		kurtosisService:   kurtosisService,
		serviceNodes:      make(map[ServiceID]ServiceNode),
		configurations:    configurations,
		testVolumeDirpath: testVolumeDirpath,
	}
}

// Gets the number of nodes in the network
func (network *ServiceNetwork) GetSize() int {
	return len(network.serviceNodes)
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
func (network *ServiceNetwork) AddService(configurationId ConfigurationID, serviceId ServiceID, dependencies map[ServiceID]bool) (*services.ServiceAvailabilityChecker, error) {
	// Maybe one day we'll make this flow from somewhere up above (e.g. make the entire network live inside a single context)
	parentCtx := context.Background()

	config, found := network.configurations[configurationId]
	if !found {
		return nil, stacktrace.NewError("No service configuration with ID '%v' has been registered", configurationId)
	}

	if _, exists := network.serviceNodes[serviceId]; exists {
		return nil, stacktrace.NewError("Service ID %s already exists in the network", serviceId)
	}

	if dependencies == nil {
		return nil, stacktrace.NewError("Dependencies map was nil; use an empty map to specify no dependencies")
	}

	// Golang maps are passed by-ref, so we do a defensive copy here so user can't change their input and mess
	// with our internal data structure
	dependencyServices := make([]services.Service, 0, len(dependencies))
	for dependencyId, _ := range dependencies  {
		dependencyNode, found := network.serviceNodes[dependencyId]
		if !found {
			return nil, stacktrace.NewError("Declared a dependency on %v but no service with this ID has been registered", dependencyId)
		}
		dependencyServices = append(dependencyServices, dependencyNode.Service)
	}


	initializer := services.NewServiceInitializer(config.initializerCore, network.testVolumeDirpath)
	service, containerId, err := initializer.CreateService(
			config.dockerImage,
			dependencyServices)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred creating service %v from configuration %v", serviceId, configurationId)
	}

	network.serviceNodes[serviceId] = ServiceNode{
		Service:     service,
		ContainerID: containerId,
	}

	availabilityChecker := services.NewServiceAvailabilityChecker(parentCtx, config.availabilityCheckerCore, service, dependencyServices)
	return availabilityChecker, nil
}

/*
Gets the node information for the service with the given service ID.
 */
func (network *ServiceNetwork) GetService(serviceId ServiceID) (ServiceNode, error) {
	node, found := network.serviceNodes[serviceId]
	if !found {
		return ServiceNode{}, stacktrace.NewError("No service with ID %v exists in the network", serviceId)
	}

	return node, nil
}

/*
Stops the container with the given service ID, and removes it from the network.
 */
func (network *ServiceNetwork) RemoveService(serviceId ServiceID, containerStopTimeoutSeconds int) error {
	nodeInfo, found := network.serviceNodes[serviceId]
	if !found {
		return stacktrace.NewError("No service with ID %v found", serviceId)
	}

	logrus.Debugf("Removing service ID %v...", serviceId)
	delete(network.serviceNodes, serviceId)

	// Make a best-effort attempt to stop the container
	err := network.kurtosisService.RemoveService(nodeInfo.ContainerID, containerStopTimeoutSeconds)
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
func (network *ServiceNetwork) RemoveAll(containerStopTimeoutSeconds int) error {
	for serviceId, _ := range network.serviceNodes {
		network.RemoveService(serviceId, containerStopTimeoutSeconds)
	}
	return nil
}
