/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package example_testsuite

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/example_impl/example_networks/single_node_example_network"
	"github.com/kurtosis-tech/kurtosis-go/lib/networks"
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type SingleNodeExampleTest struct {
	ServiceImage string
}

func (test SingleNodeExampleTest) Run(network networks.Network, context testsuite.TestContext) {
	// NOTE: We have to do this as the first line of every test because Go doesn't have generics
	castedNetwork := network.(single_node_example_network.SingleNodeExampleNetwork)

	logrus.Info("Adding the node...")
	service, err := castedNetwork.AddTheNode()
	if err != nil {
		context.Fatal(err)
	}
	logrus.Info("Successfully added the test node")

	logrus.Info("Making a query to the node...")
	serviceUrl := fmt.Sprintf("http://%v:%v", service.GetIpAddress(), service.GetPort())
	if _, err := http.Get(serviceUrl); err != nil {
		context.Fatal(stacktrace.Propagate(err, "Received an error when calling the example service endpoint"))
	}
	logrus.Info("Queried the node successfully")

	logrus.Info("Removing the node...")
	if err := castedNetwork.RemoveTheNode(); err != nil {
		context.Fatal(err)
	}
	logrus.Info("Successfully removed the test node")
}

func (test SingleNodeExampleTest) GetNetworkLoader() (networks.NetworkLoader, error) {
	return single_node_example_network.NewSingleNodeExampleNetworkLoader(test.ServiceImage), nil
}

func (test SingleNodeExampleTest) GetExecutionTimeout() time.Duration {
	return 30 * time.Second
}

func (test SingleNodeExampleTest) GetSetupBuffer() time.Duration {
	return 30 * time.Second
}

