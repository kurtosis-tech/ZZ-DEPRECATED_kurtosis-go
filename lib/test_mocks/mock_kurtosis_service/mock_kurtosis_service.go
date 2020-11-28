/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package mock_kurtosis_service

// =============== Mock Kurtosis service =====================
type MockKurtosisService struct {}

func NewMockKurtosisService() *MockKurtosisService {
	return &MockKurtosisService{}
}

func (m MockKurtosisService) AddService(
		dockerImage string,
		usedPorts map[int]bool,
		ipPlaceholder string,
		startCmdArgs []string,
		envVariables map[string]string,
		testVolumeMountLocation string) (ipAddr string, containerId string, err error) {
	return "1.2.3.4", "abcd1234", nil
}

func (m MockKurtosisService) RemoveService(containerId string, containerStopTimeoutSeconds int) error {
	return nil
}

func (m MockKurtosisService) RegisterTestExecution(testTimeoutSeconds int) error {
	return nil
}
