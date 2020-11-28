/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package services

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kurtosis-tech/kurtosis-go/lib/kurtosis_service"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type ServiceBuilder struct {
	// The core dictating initialization logic
	core DockerContainerInitializer

	// The dirpath of the service being constructed, RELATIVE TO the root of the suite execution volume!!
	serviceRelativeDirpath string

	// The handle to manipulating the test environment
	kurtosisService *kurtosis_service.KurtosisService
}

func NewServiceBuilder(core DockerContainerInitializer, serviceDirpath string, kurtosisService *kurtosis_service.KurtosisService) *ServiceBuilder {
	return &ServiceBuilder{core: core, serviceRelativeDirpath: serviceDirpath, kurtosisService: kurtosisService}
}

func (builder *ServiceBuilder) InitializeFile(fileId string, initializer func(fp *os.File) error) error {
	mountedFiles := builder.core.GetFilesToMount()
	if _, found := mountedFiles[fileId]; !found {
		return stacktrace.NewError("Attempted to initialize fileId '%v', but no such fileId was declared for the service", fileId)
	}


	return initializer()
}

func (builder *ServiceBuilder) Build() {
	logrus.Trace("Initializing files needed for service...")
	requestedFiles := initializerCore.GetFilesToMount()
	osFiles := make(map[string]*os.File)
	mountFilepaths := make(map[string]string)
	for fileId, _ := range requestedFiles {
		filename := fmt.Sprintf("%v-%v", fileId, uuid.New().String())
		hostFilepath := filepath.Join(serviceDirpath, filename)
		fp, err := os.Create(hostFilepath)
		if err != nil {
			return nil, "", stacktrace.Propagate(err, "Could not create new file for requested file ID '%v'", fileId)
		}
		defer fp.Close()
		osFiles[fileId] = fp
		mountFilepaths[fileId] = filepath.Join(mountServiceDirpath, filename)
	}
	err = initializerCore.InitializeMountedFiles(osFiles, dependencies)
	logrus.Tracef("Successfully initialized files needed for service")

	logrus.Tracef("Creating start command for service...")
	startCmdArgs, err := initializerCore.GetStartCommand(mountFilepaths, ipPlaceholder, dependencies)
	if err != nil {
		return nil, "", stacktrace.Propagate(err, "Failed to create start command.")
	}
	logrus.Tracef("Successfully created start command for service")

	logrus.Tracef("Calling to Kurtosis API to create service...")
	ipAddr, containerId, err := initializer.kurtosisService.AddService(
		dockerImage,
		usedPorts,
		ipPlaceholder,
		startCmdArgs,
		make(map[string]string),
		initializerCore.GetTestVolumeMountpoint())
	if err != nil {
		return nil, "", stacktrace.Propagate(err, "Could not add service for Docker image %v", dockerImage)
	}
	logrus.Tracef("Kurtosis API returned IP for new service: %v", ipAddr)

	logrus.Tracef("Getting service from IP...")
	service := initializer.core.GetServiceFromIp(ipAddr)
	logrus.Tracef("Successfully got service from IP")

	return service, containerId, nil
}