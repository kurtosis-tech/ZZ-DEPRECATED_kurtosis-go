/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package kurtosis_service

// =============== Mock Kurtosis service =====================
type MockKurtosisService struct {}

func NewMockKurtosisService() *MockKurtosisService {
	return &MockKurtosisService{}
}

func (m MockKurtosisService) AddService(
		serviceId string,
		dockerImage string,
		usedPorts map[string]bool,
		ipPlaceholder string,
		startCmdArgs []string,
		envVariables map[string]string,
		testVolumeMountLocation string) (ipAddr string, err error) {
	return "1.2.3.4", nil
}

func (m MockKurtosisService) RemoveService(serviceId string, containerStopTimeoutSeconds int) error {
	return nil
}

func (m MockKurtosisService) RegisterTestExecution(testTimeoutSeconds int) error {
	return nil
}
