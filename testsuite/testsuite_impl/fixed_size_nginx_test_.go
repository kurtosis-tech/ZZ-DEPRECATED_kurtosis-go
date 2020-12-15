/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package testsuite_impl

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/lib/networks"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
	"github.com/kurtosis-tech/kurtosis-go/testsuite/services_impl"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

const (
	numNodes = 5
	serviceIdPrefix = "service-"

	// Consts for starting NginX services
	timeBetweenPolls = 2 * time.Second
	numRetries = 10
)

type FixedSizeNginxTest struct {
	ServiceImage string
}

func (test FixedSizeNginxTest) Setup(context *networks.NetworkContext) (networks.Network, error) {
	containerInitializer := services_impl.NewNginxContainerInitializer(test.ServiceImage)

	availabilityCheckers := map[networks.ServiceID]services.AvailabilityChecker{}
	for i := 0; i < numNodes; i++ {
		serviceId := networks.ServiceID(serviceIdPrefix + strconv.Itoa(i))
		_, checker, err := context.AddService(serviceId, containerInitializer)
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
	return context, nil
}

func (test FixedSizeNginxTest) Run(network networks.Network, context testsuite.TestContext) {
	// NOTE: We have to do this as the first line of every test because Go doesn't have generics
	castedNetwork := network.(*networks.NetworkContext)

	for i := 0; i < numNodes; i++ {
		serviceId := networks.ServiceID(serviceIdPrefix + strconv.Itoa(i))
		logrus.Infof("Making query against service with ID '%v'...", serviceId)
		uncastedService, err := castedNetwork.GetService(serviceId)
		if err != nil {
			context.Fatal(
				stacktrace.Propagate(
					err,
					"An error occurred when getting the service interface for service with ID '%v'",
					serviceId,
				),
			)
		}
		castedService := uncastedService.(services_impl.NginxService)
		serviceUrl := fmt.Sprintf("http://%v:%v", castedService.GetIPAddress(), castedService.GetPort())
		if _, err := http.Get(serviceUrl); err != nil {
			context.Fatal(
				stacktrace.Propagate(
					err,
					"Received an error when calling the NginX service endpoint for service with ID '%v'",
					serviceId,
				),
			)
		}
		logrus.Infof("Successfully queried service with ID '%v'", serviceId)
	}
}

func (test FixedSizeNginxTest) GetExecutionTimeout() time.Duration {
	return 30 * time.Second
}

func (test FixedSizeNginxTest) GetSetupTeardownBuffer() time.Duration {
	return 30 * time.Second
}

