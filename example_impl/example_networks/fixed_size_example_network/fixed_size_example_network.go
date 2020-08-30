/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package fixed_size_example_network

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/example_impl/example_services"
	"github.com/kurtosis-tech/kurtosis-go/lib/networks"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/palantir/stacktrace"
)

const (
	vanillaConfigId    networks.ConfigurationID = "vanilla"
	vanillaDockerImage                          = "nginxdemos/hello"
	serviceIdPrefix = "service-"
)

// ======================================== NETWORK ==============================================
type FixedSizeExampleNetwork struct{
	rawNetwork *networks.ServiceNetwork
	numNodes int
}

func (network FixedSizeExampleNetwork) GetNumNodes() int {
	return network.numNodes
}

func (network *FixedSizeExampleNetwork) GetService(idInt int) (example_services.ExampleService, error) {
	if idInt < 0 || idInt >= network.numNodes {
		return nil, stacktrace.NewError("Invalid service ID '%v'", idInt)
	}
	serviceId := networks.ServiceID(fmt.Sprintf("%v%v", serviceIdPrefix, idInt))
	serviceNode, err := network.rawNetwork.GetService(serviceId)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting the service node info")
	}
	castedService := serviceNode.Service.(example_services.ExampleService)
	return castedService, nil
}

// ======================================== NETWORK LOADER ==============================================
type FixedSizeExampleNetworkLoader struct {
	numNodes int
}

func NewFixedSizeExampleNetworkLoader(numNodes int) *FixedSizeExampleNetworkLoader {
	return &FixedSizeExampleNetworkLoader{numNodes: numNodes}
}

func (loader FixedSizeExampleNetworkLoader) ConfigureNetwork(builder *networks.ServiceNetworkBuilder) error {
	builder.AddConfiguration(
		vanillaConfigId,
		vanillaDockerImage,
		example_services.ExampleServiceInitializerCore{},
		example_services.ExampleAvailabilityCheckerCore{})
	return nil
}

func (loader FixedSizeExampleNetworkLoader) InitializeNetwork(network *networks.ServiceNetwork) (map[networks.ServiceID]services.ServiceAvailabilityChecker, error) {
	availabilityCheckers := map[networks.ServiceID]services.ServiceAvailabilityChecker{}
	for i := 0; i < loader.numNodes; i++ {
		serviceId := networks.ServiceID(fmt.Sprintf("%v%v", serviceIdPrefix, i))
		checker, err := network.AddService(vanillaConfigId, serviceId, map[networks.ServiceID]bool{})
		if err != nil {
			return nil, stacktrace.Propagate(err, "An error occurred adding service with ID '%v' to the network", serviceId)
		}
		availabilityCheckers[serviceId] = *checker
	}
	return availabilityCheckers, nil
}

func (loader FixedSizeExampleNetworkLoader) WrapNetwork(network *networks.ServiceNetwork) (networks.Network, error) {
	return FixedSizeExampleNetwork{
		rawNetwork: network,
		numNodes: loader.numNodes,
	}, nil
}

