package networks

import (
	services2 "github.com/kurtosis-tech/kurtosis-go/example_impl/services"
	"github.com/kurtosis-tech/kurtosis-go/lib/networks"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
)

const (
	vanillaConfigId    networks.ConfigurationID = "vanilla"
	vanillaDockerImage                          = "nginxdemos/hello"
)

type ExampleNetworkLoader struct {}

func (e ExampleNetworkLoader) ConfigureNetwork(builder *networks.ServiceNetworkBuilder) error {
	builder.AddConfiguration(
		vanillaConfigId,
		vanillaDockerImage,
		services2.ExampleServiceInitializerCore{},
		services2.ExampleAvailabilityCheckerCore{})
	return nil
}

func (e ExampleNetworkLoader) InitializeNetwork(network *networks.ServiceNetwork) (map[networks.ServiceID]services.ServiceAvailabilityChecker, error) {
	// TODO example with some pre-test initialization
	return map[networks.ServiceID]services.ServiceAvailabilityChecker{}, nil
}

func (e ExampleNetworkLoader) WrapNetwork(network *networks.ServiceNetwork) (networks.Network, error) {
	return NewExampleNetwork(network), nil
}

