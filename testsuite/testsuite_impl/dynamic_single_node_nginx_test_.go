/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package testsuite_impl

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/lib/networks"
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
	"github.com/kurtosis-tech/kurtosis-go/testsuite/networks_impl/single_node_nginx_network"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type DynamicSingleNodeNginxTest struct {
	ServiceImage string
}

func (test DynamicSingleNodeNginxTest) Setup(context *networks.NetworkContext) (networks.Network, error) {
	return single_node_nginx_network.NewDynamicSingleNodeNginxNetwork(context, test.ServiceImage), nil
}

func (test DynamicSingleNodeNginxTest) Run(network networks.Network, context testsuite.TestContext) {
	// NOTE: We have to do this as the first line of every test because Go doesn't have generics
	castedNetwork := network.(*single_node_nginx_network.DynamicSingleNodeNginxNetwork)

	logrus.Info("Adding the node...")
	service, err := castedNetwork.AddTheNode()
	if err != nil {
		context.Fatal(err)
	}
	logrus.Info("Successfully added the test node")

	logrus.Info("Making a query to the node...")
	serviceUrl := fmt.Sprintf("http://%v:%v", service.GetIPAddress(), service.GetPort())
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

func (test DynamicSingleNodeNginxTest) GetExecutionTimeout() time.Duration {
	return 30 * time.Second
}

func (test DynamicSingleNodeNginxTest) GetSetupTeardownBuffer() time.Duration {
	return 30 * time.Second
}

