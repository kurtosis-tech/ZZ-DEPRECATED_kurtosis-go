/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package single_node_nginx_network

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/networks"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/kurtosis-tech/kurtosis-go/testsuite/services_impl"
	"github.com/palantir/stacktrace"
)

const (
	vanillaConfigId    networks.ConfigurationID = "vanilla"

	theNodeServiceId          networks.ServiceID = "the-node"
	theNodeStopTimeoutSeconds                    = 30
)

// =================================== NETWORK ===================================
type SingleNodeNginxNetwork struct{
	rawNetwork *networks.ServiceNetwork
	theNodeAdded bool
}

func NewSingleNodeNginxNetwork(rawNetwork *networks.ServiceNetwork) *SingleNodeNginxNetwork {
	return &SingleNodeNginxNetwork{
		rawNetwork: rawNetwork,
		theNodeAdded: false,
	}
}

func (network *SingleNodeNginxNetwork) AddTheNode() (services_impl.NginxService, error) {
	if network.theNodeAdded {
		return nil, stacktrace.NewError("The node is already added")
	}
	// TODO add example with dependencies
	availabilityChecker, err := network.rawNetwork.AddService(vanillaConfigId, theNodeServiceId, map[networks.ServiceID]bool{})
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred adding the node")
	}
	if err := availabilityChecker.WaitForStartup(); err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred waiting for the node to come up")
	}
	network.theNodeAdded = true

	serviceNode, err := network.rawNetwork.GetService(theNodeServiceId)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting the node's service information")
	}
	castedService := serviceNode.Service.(services_impl.NginxService)
	return castedService, nil
}

func (network *SingleNodeNginxNetwork) RemoveTheNode() error {
	if !network.theNodeAdded {
		return stacktrace.NewError("The node hasn't been added yet")
	}
	if err := network.rawNetwork.RemoveService(theNodeServiceId, theNodeStopTimeoutSeconds); err != nil {
		return stacktrace.NewError("An error occurred removing the node from the network")
	}
	network.theNodeAdded = false
	return nil
}

// =================================== NETWORK LOADER ===================================
type SingleNodeNginxNetworkLoader struct {
	serviceImage string
}

func NewSingleNodeNginxNetworkLoader(serviceImage string) *SingleNodeNginxNetworkLoader {
	return &SingleNodeNginxNetworkLoader{serviceImage: serviceImage}
}


func (loader SingleNodeNginxNetworkLoader) ConfigureNetwork(builder *networks.ServiceNetworkBuilder) error {
	builder.AddConfiguration(
		vanillaConfigId,
		loader.serviceImage,
		services_impl.NginxServiceInitializerCore{},
		services_impl.NginxAvailabilityCheckerCore{})
	return nil
}

func (loader SingleNodeNginxNetworkLoader) InitializeNetwork(network *networks.ServiceNetwork) (map[networks.ServiceID]services.AvailabilityChecker, error) {
	return map[networks.ServiceID]services.AvailabilityChecker{}, nil
}

func (loader SingleNodeNginxNetworkLoader) WrapNetwork(network *networks.ServiceNetwork) (networks.Network, error) {
	return *NewSingleNodeNginxNetwork(network), nil
}

