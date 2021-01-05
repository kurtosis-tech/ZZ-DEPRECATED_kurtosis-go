/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package kurtosis_service

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/lib/kurtosis_service/method_types"
	"github.com/palantir/stacktrace"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	kurtosisApiPort = 7443

	registrationRetryDurationSeconds = 60
	regularOperationRetryDurationSeconds = 10

	// Constants for making RPC calls to the Kurtosis API
	kurtosisServiceStruct = "KurtosisService"
	addServiceMethod = kurtosisServiceStruct + ".AddService"
	removeServiceMethod = kurtosisServiceStruct + ".RemoveService"
	repartitionMethod = kurtosisServiceStruct + ".Repartition"
	registerTestExecutionMethod = kurtosisServiceStruct + ".RegisterTestExecution"
)

// This interface provides tests with an API for performing administrative actions on the testnet, like
//  starting or stopping a service
type KurtosisService interface {
	AddService(
		serviceId 	string,
		partitionId string,
		dockerImage string,
		usedPorts map[string]bool,
		ipPlaceholder string,
		startCmdArgs []string,
		envVariables map[string]string,
		testVolumeMountLocation string) (ipAddr string, err error)

	RemoveService(serviceId string, containerStopTimeoutSeconds int) error

	Repartition(
		partitionServices map[string]map[string]bool,
		partitionConnections map[string]map[string]method_types.SerializablePartitionConnection,
		defaultConnection method_types.SerializablePartitionConnection) error

	RegisterTestExecution(testTimeoutSeconds int) error
}

type DefaultKurtosisService struct {
	ipAddr string
}

func NewDefaultKurtosisService(ipAddr string) *DefaultKurtosisService {
	return &DefaultKurtosisService{ipAddr: ipAddr}
}

func (service DefaultKurtosisService) AddService(
		serviceId string,
		partitionId string,
		dockerImage string,
		usedPorts map[string]bool,
		ipPlaceholder string,
		startCmdArgs []string,
		envVariables map[string]string,
		testVolumeMountLocation string) (ipAddr string, err error) {
	client := getConstantBackoffJsonRpcClient(service.ipAddr, regularOperationRetryDurationSeconds)
	defer client.Close()

	usedPortsList := []string{}
	for portSpecification, _ := range usedPorts {
		usedPortsList = append(usedPortsList, portSpecification)
	}
	args := method_types.AddServiceArgs{
		ServiceID: serviceId,
		PartitionID: partitionId,
		IPPlaceholder: ipPlaceholder,
		ImageName:               dockerImage,
		UsedPorts:               usedPortsList,
		StartCmd:                startCmdArgs,
		DockerEnvironmentVars:   envVariables,
		TestVolumeMountDirpath: testVolumeMountLocation,
	}
	var reply method_types.AddServiceResponse
	if err := client.Call(addServiceMethod, args, &reply); err != nil {
		return "", stacktrace.Propagate(err, "An error occurred making the call to add a service using the Kurtosis API")
	}

	return reply.IPAddress, nil
}

/*
Stops the container with the given service ID, and removes it from the network.
*/
func (service DefaultKurtosisService) RemoveService(serviceId string, containerStopTimeoutSeconds int) error {
	client := getConstantBackoffJsonRpcClient(service.ipAddr, regularOperationRetryDurationSeconds)
	defer client.Close()

	logrus.Debugf("Removing service '%v'...", serviceId)

	args := method_types.RemoveServiceArgs{
		ServiceID: serviceId,
		ContainerStopTimeoutSeconds: containerStopTimeoutSeconds,
	}

	var reply struct{}
	if err := client.Call(removeServiceMethod, args, &reply); err != nil {
		return stacktrace.Propagate(err, "An error occurred making the call to remove service '%v' using the Kurtosis API", serviceId)
	}
	logrus.Debugf("Successfully removed service '%v'", serviceId)

	return nil
}

func (service DefaultKurtosisService) Repartition(
		partitionServices map[string]map[string]bool,
		partitionConnections map[string]map[string]method_types.SerializablePartitionConnection,
		defaultConnection method_types.SerializablePartitionConnection) error {
	client := getConstantBackoffJsonRpcClient(service.ipAddr, regularOperationRetryDurationSeconds)
	defer client.Close()

	logrus.Debugf("Repartitioning test network with the following args:")
	logrus.Debugf("New partition services: %v", partitionServices)
	logrus.Debugf("New partition connections: %v", partitionConnections)
	logrus.Debugf("New default connection: %v", defaultConnection)

	args := method_types.RepartitionArgs{
		PartitionServices:    partitionServices,
		PartitionConnections: partitionConnections,
		DefaultConnection:    defaultConnection,
	}

	var reply struct{}
	if err := client.Call(repartitionMethod, args, &reply); err != nil {
		return stacktrace.Propagate(err, "An error occurred making the call to repartition the test network using the Kurtosis API")
	}
	logrus.Debugf("Successfully repartitioned the test network")
	return nil
}

func (service DefaultKurtosisService) RegisterTestExecution(testTimeoutSeconds int) error {
	client := getConstantBackoffJsonRpcClient(service.ipAddr, registrationRetryDurationSeconds)
	defer client.Close()

	logrus.Debugf("Registering a test execution with a timeout of %v seconds...", testTimeoutSeconds)

	args := method_types.RegisterTestExecutionArgs{TestTimeoutSeconds: testTimeoutSeconds}

	var reply struct{}
	if err := client.Call(registerTestExecutionMethod, args, &reply); err != nil {
		return stacktrace.Propagate(err, "An error occurred making the call to register a test execution using the Kurtosis API")
	}
	logrus.Debugf("Successfully registered a test execution with timeout of %v seconds", testTimeoutSeconds)

	return nil

}

// ================================= Private helper function ============================================
func getConstantBackoffJsonRpcClient(ipAddr string, retryDurationSeconds int) *jsonrpc2.Client {
	kurtosisUrl := fmt.Sprintf("http://%v:%v", ipAddr, kurtosisApiPort)

	// We intentionally don't use a retrying client because if an error occurs on the Kurtosis API, a
	//  retrying client would silently retry. This can lead to distracting errors like "cannot add
	//  service ID X; it already exists!" when the real problem was with the first call to add the service.
	return jsonrpc2.NewCustomHTTPClient(kurtosisUrl, &http.Client{})
}
