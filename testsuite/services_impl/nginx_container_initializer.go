/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package services_impl

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"os"
)

const (
	testVolumeMountpoint = "/test-volume"
)

type NginxContainerInitializer struct{
	dockerImage string
}

func NewNginxContainerInitializer(dockerImage string) *NginxContainerInitializer {
	return &NginxContainerInitializer{dockerImage: dockerImage}
}

func (initializer NginxContainerInitializer) GetDockerImage() string {
	return initializer.dockerImage
}

func (initializer NginxContainerInitializer) GetUsedPorts() map[int]bool {
	return map[int]bool{
		nginxServicePort: true,
	}
}

func (initializer NginxContainerInitializer) GetServiceFromIp(ipAddr string) services.Service {
	return NginxService{IPAddr: ipAddr}
}

func (initializer NginxContainerInitializer) GetFilesToMount() map[string]bool {
	// TODO give an example of mounting files
	return map[string]bool{}
}

func (initializer NginxContainerInitializer) InitializeMountedFiles(mountedFiles map[string]*os.File) error {
	// TODO give example of mounting files
	return nil
}

func (initializer NginxContainerInitializer) GetTestVolumeMountpoint() string {
	return testVolumeMountpoint
}

func (initializer NginxContainerInitializer) GetStartCommand(mountedFileFilepaths map[string]string, ipPlaceholder string) ([]string, error) {
	// If there was a specific start command that we wanted Docker to run, we'd return the string array here. By
	//	returning nil, we tell Kurtosis to run the image with whatever CMD or ENTRYPOINT is specified in the
	//	Dockerfile will be used instead. This prevents the Kurtosis code from needing to know specifics about
	//	the Docker image (initializer.g. what filepath the binary to run is located at)
	return nil, nil
}
