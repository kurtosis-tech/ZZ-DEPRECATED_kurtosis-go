package services

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"os"
)

const (
	exampleServicePort = 80
	exampleServiceTestVolumeMountpoint = "/shared"
)

type ExampleServiceInitializerCore struct{}

func (e ExampleServiceInitializerCore) GetUsedPorts() map[int]bool {
	return map[int]bool{
		exampleServicePort: true,
	}
}

func (e ExampleServiceInitializerCore) GetServiceFromIp(ipAddr string) services.Service {
	return Socket{
		IPAddr: ipAddr,
		Port:   exampleServicePort,
	}
}

func (e ExampleServiceInitializerCore) GetFilesToMount() map[string]bool {
	// TODO give an example of mounting files
	return map[string]bool{}
}

func (e ExampleServiceInitializerCore) InitializeMountedFiles(mountedFiles map[string]*os.File, dependencies []services.Service) error {
	// TODO give example of mounting files
	return nil
}

func (e ExampleServiceInitializerCore) GetTestVolumeMountpoint() string {
	return exampleServiceTestVolumeMountpoint
}

func (e ExampleServiceInitializerCore) GetStartCommand(mountedFileFilepaths map[string]string, ipPlaceholder string, dependencies []services.Service) ([]string, error) {
	// TODO Explain why this is nil, or maybe make entrypoints/environment-variable-flavored launches more obvious?
	return nil, nil
}
