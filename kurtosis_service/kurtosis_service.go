package kurtosis_service

import (
	"context"
	"fmt"
	"github.com/palantir/stacktrace"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	addServiceMethod = "KurtosisAPI.AddService"
	removeServiceMethod = "KurtosisAPI.RemoveService"
)

type KurtosisService struct {
	ipAddr string
	port string // TODO change this type
}

func NewKurtosisService(ipAddr string, port string) *KurtosisService {
	return &KurtosisService{ipAddr: ipAddr, port: port}
}

func (service KurtosisService) AddService(
		dockerImage string,
		// TODO change type of this
		usedPorts map[string]bool,
		ipPlaceholder string,
		startCmdArgs []string,
		envVariables map[string]string,
		testVolumeMountLocation string) (string, error) {

	// TODO reuse clients?
	client := jsonrpc2.NewHTTPClient(fmt.Sprintf("%v:%v", service.ipAddr, service.port))
	defer client.Close()

	usedPortsList := []string{}
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
		return "", stacktrace.Propagate(err, "An error occurred making the call to add a service using the Kurtosis API")
	}

	return reply.IPAddress, nil
}

/*
Stops the container with the given service ID, and removes it from the network.
*/
func (api KurtosisService) RemoveService(serviceId ServiceID) error {
	// TODO reuse clients?
	client := jsonrpc2.NewHTTPClient(fmt.Sprintf("%v:%v", service.ipAddr, service.port))
	defer client.Close()

	logrus.Debugf("Removing service with ID %v...", serviceId)

	args := RemoveServiceArgs{
		ServiceID: serviceId,
	}

	var reply struct{}
	if err := client.Call(addServiceMethod, args, &reply); err != nil {
		return stacktrace.Propagate(err, "An error occurred making the call to remove a service using the Kurtosis API")
	}
	logrus.Debugf("Successfully removed service with ID %v", serviceId)

	return nil
}
