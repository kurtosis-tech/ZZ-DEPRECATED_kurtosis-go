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

type SingleNodeExampleTest struct {}

func (e SingleNodeExampleTest) Run(network networks.Network, context testsuite.TestContext) {
	// NOTE: We have to do this as the first line of every test because Go doesn't have generics
	castedNetwork := network.(single_node_example_network.SingleNodeExampleNetwork)

	logrus.Info("Adding the node...")
	service, err := castedNetwork.AddTheNode()
	if err != nil {
		context.Fatal(err)
	}
	logrus.Info("Successfully added the test node")

	logrus.Info("Making a query to the node...")
	socket := service.GetHelloWorldSocket()
	serviceUrl := fmt.Sprintf("http://%v:%v", socket.IPAddr, socket.Port)
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

func (e SingleNodeExampleTest) GetNetworkLoader() (networks.NetworkLoader, error) {
	return single_node_example_network.SingleNodeExampleNetworkLoader{}, nil
}

func (e SingleNodeExampleTest) GetExecutionTimeout() time.Duration {
	return 30 * time.Second
}

func (e SingleNodeExampleTest) GetSetupBuffer() time.Duration {
	return 30 * time.Second
}

