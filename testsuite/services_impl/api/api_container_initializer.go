/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package api

import (
	"encoding/json"
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/kurtosis-tech/kurtosis-go/testsuite/services_impl/datastore"
	"github.com/palantir/stacktrace"
	"os"
)

const (
	port = 2434

	configFileKey = "config-file"

	testVolumeMountpoint = "/test-volume"
)

type config struct {
	datastoreIp string	`json:"datastoreIp"`
	datastorePort int	`json:"datastorePort"`
}

type ApiContainerInitializer struct {
	dockerImage string
	datastore datastore.DatastoreService
}

func NewApiContainerInitializer(dockerImage string, datastore datastore.DatastoreService) *ApiContainerInitializer {
	return &ApiContainerInitializer{dockerImage: dockerImage, datastore: datastore}
}

func (initializer ApiContainerInitializer) GetDockerImage() string {
	return initializer.dockerImage
}

func (initializer ApiContainerInitializer) GetUsedPorts() map[string]bool {
	return map[string]bool{
		fmt.Sprintf("%v/tcp", port): true,
	}
}

func (initializer ApiContainerInitializer) GetServiceFromIp(ipAddr string) services.Service {
	return NewApiService(ipAddr, port)
}

func (initializer ApiContainerInitializer) GetFilesToMount() map[string]bool {
	return map[string]bool{
		configFileKey: true,
	}
}

func (initializer ApiContainerInitializer) InitializeMountedFiles(mountedFiles map[string]*os.File) error {
	configObj := config{
		datastoreIp:   initializer.datastore.GetIPAddress(),
		datastorePort: initializer.datastore.GetPort(),
	}
	configBytes, err := json.Marshal(&configObj)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred serializing the config to JSON")
	}

	configFp := mountedFiles[configFileKey]
	if _, err := configFp.Write(configBytes); err != nil {
		return stacktrace.Propagate(err, "An error occurred writing the serialized config JSON to file")
	}

	return nil
}

func (initializer ApiContainerInitializer) GetTestVolumeMountpoint() string {
	return testVolumeMountpoint
}

func (initializer ApiContainerInitializer) GetStartCommand(mountedFileFilepaths map[string]string, ipPlaceholder string) ([]string, error) {
	// The Dockerfile for the service has a start command baked in, so we don't need to specify one
	return nil, nil
}

