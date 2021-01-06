/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package nginx_static

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"os"
	"strconv"
)

type NginxStaticContainerInitializer struct {
	artifactUrl string
}

func NewNginxStaticContainerInitializer(artifactUrl string) *NginxStaticContainerInitializer {
	return &NginxStaticContainerInitializer{artifactUrl: artifactUrl}
}

func (s NginxStaticContainerInitializer) GetDockerImage() string {
	return dockerImage
}

func (s NginxStaticContainerInitializer) GetUsedPorts() map[string]bool {
	return map[string]bool{
		strconv.Itoa(listenPort): true,
	}
}

func (s NginxStaticContainerInitializer) GetService(serviceId services.ServiceID, ipAddr string) services.Service {
	return &NginxStaticService{
		serviceId: serviceId,
		ipAddr:    ipAddr,
	}
}

func (s NginxStaticContainerInitializer) GetFilesToMount() map[string]bool {
	// No generated files to mount
	return map[string]bool{}
}

func (s NginxStaticContainerInitializer) InitializeMountedFiles(mountedFiles map[string]*os.File) error {
	// No generated files to initialize
	return nil
}

func (s NginxStaticContainerInitializer) GetTestVolumeMountpoint() string {
	return "/test-volume"
}

func (s NginxStaticContainerInitializer) GetStartCommand(mountedFileFilepaths map[string]string, ipPlaceholder string) ([]string, error) {
	// Don't specify an explicit start command - default to using the baked-in command
	return nil, nil
}
