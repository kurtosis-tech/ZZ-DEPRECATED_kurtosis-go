/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package mock_service

const (
	MockServicePort = 1000
)

// Mock service, for testing purposes only
type MockService struct {
	ipAddr string

	// For testing, the service will report as available on the Nth call to IsAvailable
	becomesAvailableOnCheck int

	// Number of calls to IsAvailable that have happened
	callsToIsAvailable int
}

func NewMockService(ipAddr string, becomesAvailableOnCheck int) *MockService {
	return &MockService{
		ipAddr:                  ipAddr,
		becomesAvailableOnCheck: becomesAvailableOnCheck,
		callsToIsAvailable:      0,
	}
}


func (m MockService) GetIPAddress() string {
	return m.ipAddr
}

func (m *MockService) IsAvailable() bool {
	m.callsToIsAvailable++
	return m.callsToIsAvailable >= m.becomesAvailableOnCheck
}

