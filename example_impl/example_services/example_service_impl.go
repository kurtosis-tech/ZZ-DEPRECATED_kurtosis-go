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

func (e ExampleServiceImpl) GetHelloWorldSocket() Socket {
	return Socket{
		IPAddr: e.IPAddr,
		Port: exampleServicePort,
	}
}
