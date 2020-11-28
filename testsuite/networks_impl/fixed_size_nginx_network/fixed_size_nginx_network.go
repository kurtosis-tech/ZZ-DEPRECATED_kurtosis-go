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
	"time"
)

const (
	serviceIdPrefix = "service-"

	// Consts for starting NginX services
	timeBetweenPolls = 2 * time.Second
	numRetries = 10
)

// ======================================== NETWORK ==============================================
type NginxNetwork struct{
	// Network context
	networkCtx *networks.NetworkContext

	// Number of nodes in the network
	numNodes int
}

func NewFixedSizeNginxNetwork(networkCtx *networks.NetworkContext, dockerImage string, numNodes int) (*NginxNetwork, error) {
	// Meta: It's not great that we do actual logic here in the constructor; if this becomes problematic
	//  then we can create a Builder to separate the configuration & instantiation
	availabilityCheckers := map[networks.ServiceID]services.AvailabilityChecker{}
	for i := 0; i < numNodes; i++ {
		serviceId := networks.ServiceID(fmt.Sprintf("%v%v", serviceIdPrefix, i))
		_, checker, err := networkCtx.AddService(serviceId, services_impl.NewNginxContainerInitializer(dockerImage))
		if err != nil {
			return nil, stacktrace.Propagate(err, "An error occurred adding service with ID '%v' to the network", serviceId)
		}
		availabilityCheckers[serviceId] = checker
	}

	// Now that all services are started, wait for them to come up
	for serviceId, checker := range availabilityCheckers {
		if err := checker.WaitForStartup(timeBetweenPolls, numRetries); err != nil {
			return nil, stacktrace.Propagate(err, "An error occurred waiting for service ID '%v' to come up", serviceId)
		}
	}
	return &NginxNetwork{networkCtx: networkCtx, numNodes: numNodes}, nil
}

func (network NginxNetwork) GetNumNodes() int {
	return network.numNodes
}

func (network *NginxNetwork) GetService(idInt int) (services_impl.NginxService, error) {
	if idInt < 0 || idInt >= network.numNodes {
		return services_impl.NginxService{}, stacktrace.NewError("Invalid service ID '%v'", idInt)
	}
	serviceId := networks.ServiceID(fmt.Sprintf("%v%v", serviceIdPrefix, idInt))
	service, err := network.networkCtx.GetService(serviceId)
	if err != nil {
		return services_impl.NginxService{}, stacktrace.Propagate(err, "An error occurred getting the service node info")
	}
	castedService := service.(services_impl.NginxService)
	return castedService, nil
}
