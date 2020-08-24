package services

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kurtosis-tech/kurtosis-go/lib/kurtosis_service"
	"github.com/palantir/stacktrace"
	"net"
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
func NewServiceInitializer(core ServiceInitializerCore, testVolumeDirpath string) *ServiceInitializer {
	return &ServiceInitializer{
		core: core,
		testVolumeDirpath: testVolumeDirpath,
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

	serviceDirname := fmt.Sprintf("service-%v", uuid.New().String())
	// TODO figure out a better way to do this; the testsuite might collide with the Kurtosis API!!!
	serviceDirpath := filepath.Join(initializer.testVolumeDirpath, serviceDirname)
	err := os.Mkdir(serviceDirpath, os.ModeDir)
	if err != nil {
		return nil, "", stacktrace.Propagate(err, "An error occurred creating the new service's directory in the volume at filepath '%v'", serviceDirpath)
	}
	mountServiceDirpath := filepath.Join(initializerCore.GetTestVolumeMountpoint(), serviceDirname)

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
	startCmdArgs, err := initializerCore.GetStartCommand(mountFilepaths, ipPlaceholder, dependencies)
	if err != nil {
		return nil, "", stacktrace.Propagate(err, "Failed to create start command.")
	}

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
	return initializer.core.GetServiceFromIp(ipAddr), containerId, nil
}

/*
Calls down to the initializer core to get an instance of the user-defined interface that is used for interacting with
	the user's service. The core will do the instantiation of the actual interface implementation.
 */
func (initializer ServiceInitializer) GetServiceFromIp(ipAddr net.IP) Service {
	return initializer.core.GetServiceFromIp(ipAddr.String())
}
