package networks

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/networks"
	"github.com/palantir/stacktrace"
)

const (
	theNodeServiceId          networks.ServiceID = "the-node"
	theNodeStopTimeoutSeconds                    = 30
)

type ExampleNetwork struct{
	rawNetwork *networks.ServiceNetwork
	theNodeAdded bool
}

func NewExampleNetwork(rawNetwork *networks.ServiceNetwork) *ExampleNetwork {
	return &ExampleNetwork{
		rawNetwork: rawNetwork,
		theNodeAdded: false,
	}
}

func (network *ExampleNetwork) AddTheNode() error {
	if network.theNodeAdded {
		return stacktrace.NewError("The node is already added")
	}
	// TODO add example with dependencies
	network.rawNetwork.AddService(vanillaConfigId, theNodeServiceId, map[networks.ServiceID]bool{})
	return nil
}

func (network *ExampleNetwork) RemoveTheNode() error {
	if !network.theNodeAdded {
		return stacktrace.NewError("The node hasn't been added yet")
	}
	network.rawNetwork.RemoveService(theNodeServiceId, theNodeStopTimeoutSeconds)
	return nil
}
