/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package networks

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/kurtosis_service"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/palantir/stacktrace"
)

// Identifier used for service configurations
type ConfigurationID string

/*
A builder for configuring & constructing a test ServiceNetwork.
 */
type ServiceNetworkBuilder struct {
	// The Kurtosis service that will be used for manipulating the Docker engine during the test
	kurtosisService *kurtosis_service.KurtosisService

	// Mapping of configuration ID -> factories used to construct new nodes
	configurations map[ConfigurationID]serviceConfig

	// Directory path where the test Docker volume is mounted on the test suite image
	testVolumeDirpath string
}

/*
Creates a new builder for configuring a ServiceNetwork.

Args:
	kurtosisService: Docker manager that will be used to manipulate the Docker engine when adding services
	testVolumeDirpath: The dirpath where the test volume is mounted on the controller (which is where this code
		will be executing)
 */
func NewServiceNetworkBuilder(kurtosisService *kurtosis_service.KurtosisService, testVolumeDirpath string) *ServiceNetworkBuilder {
	return &ServiceNetworkBuilder{
		kurtosisService:   kurtosisService,
		configurations:    map[ConfigurationID]serviceConfig{},
		testVolumeDirpath: testVolumeDirpath,
	}
}

/*
Defines a new service configuration to the network that can later be used to launch Docker containers

Args:
	configurationId: The ID by which this configuration will be referenced later
	dockerImage: The Docker image that containers launched with this configuration will run with
	initializerCore: The user-defined logic for how to launch the Docker container
	availabilityCheckerCore: The user-defined logic for how to report services launched with this configuration
		as available
 */
func (builder *ServiceNetworkBuilder) AddConfiguration(
			configurationId ConfigurationID,
			dockerImage string,
			initializerCore services.ServiceInitializerCore,
			availabilityCheckerCore services.ServiceAvailabilityCheckerCore) error {
	if _, found := builder.configurations[configurationId]; found {
		return stacktrace.NewError("Configuration ID %v is already registered", configurationId)
	}

	serviceConfig := serviceConfig{
		dockerImage: dockerImage,
		availabilityCheckerCore: availabilityCheckerCore,
		initializerCore:         initializerCore,
	}
	builder.configurations[configurationId] = serviceConfig
	return nil
}

/*
Constructs a ServiceNetwork with the configurations that were defined for this builder
 */
func (builder ServiceNetworkBuilder) Build() *ServiceNetwork {
	// Defensive copy, so user calling functions on the builder after building won't affect the
	// state of the object we already built
	configurationsCopy := make(map[ConfigurationID]serviceConfig)
	for configurationId, config := range builder.configurations {
		configurationsCopy[configurationId] = config
	}
	return NewServiceNetwork(
		builder.kurtosisService,
		configurationsCopy,
		builder.testVolumeDirpath)
}
