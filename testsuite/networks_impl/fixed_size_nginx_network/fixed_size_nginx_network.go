/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package fixed_size_nginx_network

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/lib/networks"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/kurtosis-tech/kurtosis-go/testsuite/services_impl"
	"github.com/palantir/stacktrace"
)

const (
	vanillaConfigId    networks.ConfigurationID = "vanilla"
	serviceIdPrefix = "service-"
)

// ======================================== NETWORK ==============================================
type FixedSizeNginxNetwork struct{
	rawNetwork *networks.ServiceNetwork
	numNodes int
}

func (network FixedSizeNginxNetwork) GetNumNodes() int {
	return network.numNodes
}

func (network *FixedSizeNginxNetwork) GetService(idInt int) (services_impl.NginxService, error) {
	if idInt < 0 || idInt >= network.numNodes {
		return nil, stacktrace.NewError("Invalid service ID '%v'", idInt)
	}
	serviceId := networks.ServiceID(fmt.Sprintf("%v%v", serviceIdPrefix, idInt))
	serviceNode, err := network.rawNetwork.GetService(serviceId)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting the service node info")
	}
	castedService := serviceNode.Service.(services_impl.NginxService)
	return castedService, nil
}

// ======================================== NETWORK LOADER ==============================================
type FixedSizeNginxNetworkLoader struct {
	numNodes int
	serviceImage string
}

func NewFixedSizeNginxNetworkLoader(numNodes int, serviceImage string) *FixedSizeNginxNetworkLoader {
	return &FixedSizeNginxNetworkLoader{
		numNodes: numNodes,
		serviceImage: serviceImage,
	}
}

func (loader FixedSizeNginxNetworkLoader) ConfigureNetwork(builder *networks.ServiceNetworkBuilder) error {
	builder.AddConfiguration(
		vanillaConfigId,
		loader.serviceImage,
		services_impl.NginxServiceInitializerCore{},
		services_impl.NginxAvailabilityCheckerCore{})
	return nil
}

func (loader FixedSizeNginxNetworkLoader) InitializeNetwork(network *networks.ServiceNetwork) (map[networks.ServiceID]services.ServiceAvailabilityChecker, error) {
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

func (loader FixedSizeNginxNetworkLoader) WrapNetwork(network *networks.ServiceNetwork) (networks.Network, error) {
	return FixedSizeNginxNetwork{
		rawNetwork: network,
		numNodes: loader.numNodes,
	}, nil
}

