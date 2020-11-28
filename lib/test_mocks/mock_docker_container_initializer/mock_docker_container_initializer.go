/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package mock_docker_container_initializer

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/kurtosis-tech/kurtosis-go/lib/test_mocks/mock_service"
	"os"
)

type MockDockerContainerInitializer struct{}

func NewMockDockerContainerInitializer() *MockDockerContainerInitializer {
	return &MockDockerContainerInitializer{}
}

func (m MockDockerContainerInitializer) GetDockerImage() string {
	return "some-image"
}

func (m MockDockerContainerInitializer) GetUsedPorts() map[int]bool {
	return map[int]bool{
		mock_service.MockServicePort: true,
	}
}

func (m MockDockerContainerInitializer) GetServiceFromIp(ipAddr string) services.Service {
	return mock_service.NewMockService(ipAddr, true)
}

func (m MockDockerContainerInitializer) GetFilesToMount() map[string]bool {
	return map[string]bool{}
}

func (m MockDockerContainerInitializer) InitializeMountedFiles(mountedFiles map[string]*os.File) error {
	return nil
}

func (m MockDockerContainerInitializer) GetTestVolumeMountpoint() string {
	return "/test-volume"
}

func (m MockDockerContainerInitializer) GetStartCommand(mountedFileFilepaths map[string]string, ipPlaceholder string) ([]string, error) {
	return []string{
		"some-binary",
		"--some-flag",
	}, nil
}

