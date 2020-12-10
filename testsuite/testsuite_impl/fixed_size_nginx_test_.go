/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package testsuite_impl

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/lib/networks"
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
	"github.com/kurtosis-tech/kurtosis-go/testsuite/networks_impl/fixed_size_nginx_network"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	numNodes = 5
)

type FixedSizeNginxTest struct {
	ServiceImage string
}

func (test FixedSizeNginxTest) Setup(context *networks.NetworkContext) (networks.Network, error) {
	network, err := fixed_size_nginx_network.NewFixedSizeNginxNetwork(context, test.ServiceImage, numNodes)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred setting up the test network")
	}
	return network, nil
}

func (test FixedSizeNginxTest) Run(network networks.Network, context testsuite.TestContext) {
	// NOTE: We have to do this as the first line of every test because Go doesn't have generics
	castedNetwork := network.(*fixed_size_nginx_network.FixedSizeNginxNetwork)

	for i := 0; i < castedNetwork.GetNumNodes(); i++ {
		logrus.Infof("Making query against node #%v...", i)
		service, err := castedNetwork.GetService(i)
		if err != nil {
			context.Fatal(stacktrace.Propagate(err, "An error occurred when getting the service interface for node #%v", i))
		}
		serviceUrl := fmt.Sprintf("http://%v:%v", service.GetIPAddress(), service.GetPort())
		if _, err := http.Get(serviceUrl); err != nil {
			context.Fatal(stacktrace.Propagate(err, "Received an error when calling the example service endpoint for node #%v", i))
		}
		logrus.Infof("Successfully queried node #%v", i)
	}
}

func (test FixedSizeNginxTest) GetExecutionTimeout() time.Duration {
	return 30 * time.Second
}

func (test FixedSizeNginxTest) GetSetupTeardownBuffer() time.Duration {
	return 30 * time.Second
}

