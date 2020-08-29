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

const (
	ipPlaceholder = "KURTOSISSERVICEIP"
)

// TODO We MIGHT be able to remove this struct entirely
/*
A struct that wraps a user-defined ServiceInitializerCore, which will instruct the initializer how to launch a new instance
	of the user's service.
 */
type ServiceInitializer struct {
	// The user-defined instructions for how to initialize their service
	core ServiceInitializerCore

	// The location where the test volume is mounted *on the test suite container*
	testVolumeDirpath string

	// The handle to manipulating the test environment
	kurtosisService *kurtosis_service.KurtosisService
}

/*
Creates a new service initializer that will initialize services using the user-defined core.

Args:
	core: The user-defined logic for instantiating their particular service
	testVolumeDirpath: The dirpath where the test Docker volume is mounted on the test suite Docker container
 */
func NewServiceInitializer(
		core ServiceInitializerCore,
		testVolumeDirpath string,
		kurtosisService *kurtosis_service.KurtosisService) *ServiceInitializer {
	return &ServiceInitializer{
		core: core,
		testVolumeDirpath: testVolumeDirpath,
		kurtosisService: kurtosisService,
	}
}

// If Go had generics, this would be genericized so that the arg type = return type
/*
Creates a service with the given parameters

Args:
	dockerImage: The name of the Docker image that the new service will be started with
	ipPlaceholder: Since the user won't know the IP address of the service being created in advance, this is the
		placeholder string that will be used instead (and which will be swapped with the actual IP before service
		launch)
	dependencies: The services that the service-to-be-started depends on

Returns:
	Service: The interface which should be used to access the newly-created service (which, because Go doesn't have generics,
		will need to be casted to the appropriate type)
	string: The ID of the service as returned by the Kurtosis API
 */
func (initializer ServiceInitializer) CreateService(
			dockerImage string,
			dependencies []Service) (Service, string, error) {
	initializerCore := initializer.core
	usedPorts := initializerCore.GetUsedPorts()

	logrus.Trace("Creating directory within test volume for service...")
	serviceDirname := fmt.Sprintf("service-%v", uuid.New().String())
	// TODO figure out a better way to do this; the testsuite might collide with the Kurtosis API!!!
	serviceDirpath := filepath.Join(initializer.testVolumeDirpath, serviceDirname)
	err := os.Mkdir(serviceDirpath, os.ModeDir)
	if err != nil {
		return nil, "", stacktrace.Propagate(err, "An error occurred creating the new service's directory in the volume at filepath '%v'", serviceDirpath)
	}
	mountServiceDirpath := filepath.Join(initializerCore.GetTestVolumeMountpoint(), serviceDirname)
	logrus.Trace("Successfully created directory within test volume for service")

	logrus.Trace("Initializing files needed for service...")
	requestedFiles := initializerCore.GetFilesToMount()
	osFiles := make(map[string]*os.File)
	mountFilepaths := make(map[string]string)
	for fileId, _ := range requestedFiles {
		filename := uuid.New().String()
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
