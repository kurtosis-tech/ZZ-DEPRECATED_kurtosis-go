package kurtosis_service

import (
	"fmt"
	"github.com/palantir/stacktrace"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/sirupsen/logrus"
)

const (
	kurtosisApiPort = 7443
	addServiceMethod = "KurtosisAPI.AddService"
	removeServiceMethod = "KurtosisAPI.RemoveService"
)

type KurtosisService struct {
	ipAddr string
}

func NewKurtosisService(ipAddr string) *KurtosisService {
	return &KurtosisService{ipAddr: ipAddr}
}

/*
Calls the Kurtosis API container to add a service to the network
 */
func (service KurtosisService) AddService(
		dockerImage string,
		// TODO change type of this to be an actual Port type
		usedPorts map[int]bool,
		ipPlaceholder string,
		startCmdArgs []string,
		envVariables map[string]string,
		testVolumeMountLocation string) (string, string, error) {
	client := jsonrpc2.NewHTTPClient(fmt.Sprintf("%v:%v", service.ipAddr, kurtosisApiPort))
	defer client.Close()

	// TODO allow non-TCP protocols
	usedPortsList := []int{}
	for port, _ := range usedPorts {
		usedPortsList = append(usedPortsList, port)
	}
	args := AddServiceArgs{
		IPPlaceholder: ipPlaceholder,
		ImageName:               dockerImage,
		UsedPorts:               usedPortsList,
		StartCmd:                startCmdArgs,
		DockerEnvironmentVars:   envVariables,
		TestVolumeMountFilepath: testVolumeMountLocation,
	}
	var reply AddServiceResponse
	if err := client.Call(addServiceMethod, args, &reply); err != nil {
		return "", "", stacktrace.Propagate(err, "An error occurred making the call to add a service using the Kurtosis API")
	}

	return reply.IPAddress, reply.ContainerID, nil
}

/*
Stops the container with the given service ID, and removes it from the network.
*/
func (service KurtosisService) RemoveService(containerId string, containerStopTimeoutSeconds int) error {
	client := jsonrpc2.NewHTTPClient(fmt.Sprintf("%v:%v", service.ipAddr, kurtosisApiPort))
	defer client.Close()

	logrus.Debugf("Removing service with container ID %v...", containerId)

	args := RemoveServiceArgs{
		ContainerID: containerId,
		ContainerStopTimeoutSeconds: containerStopTimeoutSeconds,
	}

	var reply struct{}
	if err := client.Call(addServiceMethod, args, &reply); err != nil {
		return stacktrace.Propagate(err, "An error occurred making the call to remove a service using the Kurtosis API")
	}
	logrus.Debugf("Successfully removed service with container ID %v", containerId)

	return nil
}
