package testsuite

import (
	networks2 "github.com/kurtosis-tech/kurtosis-go/example_impl/networks"
	"github.com/kurtosis-tech/kurtosis-go/lib/networks"
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
	"github.com/sirupsen/logrus"
	"time"
)

type ExampleTest struct {}

func (e ExampleTest) Run(network networks.Network, context testsuite.TestContext) {
	logrus.Info("Adding the node...")
	castedNetwork := network.(networks2.ExampleNetwork)
	if err := castedNetwork.AddTheNode(); err != nil {
		context.Fatal(err)
	}
	logrus.Info("Successfully added the test node")

	logrus.Info("Removing the node...")
	if err := castedNetwork.RemoveTheNode(); err != nil {
		context.Fatal(err)
	}
	logrus.Info("Successfully removed the test node")
}

func (e ExampleTest) GetNetworkLoader() (networks.NetworkLoader, error) {
	return networks2.ExampleNetworkLoader{}, nil
}

func (e ExampleTest) GetExecutionTimeout() time.Duration {
	return 30 * time.Second
}

func (e ExampleTest) GetSetupBuffer() time.Duration {
	return 10 * time.Second
}

