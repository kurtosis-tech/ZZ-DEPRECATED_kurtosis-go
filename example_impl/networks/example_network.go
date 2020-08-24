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
	availabilityChecker, err := network.rawNetwork.AddService(vanillaConfigId, theNodeServiceId, map[networks.ServiceID]bool{})
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred adding the node")
	}
	if err := availabilityChecker.WaitForStartup(); err != nil {
		return stacktrace.Propagate(err, "An error occurred waiting for the node to come up")
	}
	network.theNodeAdded = true
	return nil
}

func (network *ExampleNetwork) RemoveTheNode() error {
	if !network.theNodeAdded {
		return stacktrace.NewError("The node hasn't been added yet")
	}
	if err := network.rawNetwork.RemoveService(theNodeServiceId, theNodeStopTimeoutSeconds); err != nil {
		return stacktrace.NewError("An error occurred removing the node from the network")
	}
	network.theNodeAdded = false
	return nil
}
