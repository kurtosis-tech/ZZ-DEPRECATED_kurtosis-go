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
	"time"
)

const (
	nodeServiceId networks.ServiceID = "the-node"

	// Consts for starting & stopping the node
	timeBetweenPolls   = 2 * time.Second
	numRetries         = 10
	stopTimeoutSeconds = 30
)

// =================================== NETWORK ===================================
type DynamicSingleNodeNginxNetwork struct{
	networkCtx        *networks.NetworkContext
	dockerInitializer services.DockerContainerInitializer
	nodeAdded         bool
}

func NewDynamicSingleNodeNginxNetwork(networkCtx *networks.NetworkContext, dockerImage string) *DynamicSingleNodeNginxNetwork {
	return &DynamicSingleNodeNginxNetwork{
		networkCtx:        networkCtx,
		dockerInitializer: services_impl.NewNginxContainerInitializer(dockerImage),
		nodeAdded:         false,
	}
}

func (network *DynamicSingleNodeNginxNetwork) AddTheNode() (services_impl.NginxService, error) {
	if network.nodeAdded {
		return services_impl.NginxService{}, stacktrace.NewError("The node is already added")
	}
	// TODO add example with dependencies
	service, availabilityChecker, err := network.networkCtx.AddService(nodeServiceId, network.dockerInitializer)
	if err != nil {
		return services_impl.NginxService{}, stacktrace.Propagate(err, "An error occurred adding the node")
	}
	if err := availabilityChecker.WaitForStartup(timeBetweenPolls, numRetries); err != nil {
		return services_impl.NginxService{}, stacktrace.Propagate(err, "An error occurred waiting for the node to come up")
	}
	network.nodeAdded = true
	castedService := service.(services_impl.NginxService)
	return castedService, nil
}

func (network *DynamicSingleNodeNginxNetwork) RemoveTheNode() error {
	if !network.nodeAdded {
		return stacktrace.NewError("The node hasn't been added yet")
	}
	if err := network.networkCtx.RemoveService(nodeServiceId, stopTimeoutSeconds); err != nil {
		return stacktrace.NewError("An error occurred removing the node from the network")
	}
	network.nodeAdded = false
	return nil
}
