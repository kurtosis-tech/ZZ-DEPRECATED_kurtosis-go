/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package example_services

const (
	exampleServicePort = 80
)
type ExampleServiceImpl struct{
	IPAddr string
}

func (e ExampleServiceImpl) GetIpAddress() string {
	return e.IPAddr
}

func (e ExampleServiceImpl) GetPort() int {
	return exampleServicePort
}
