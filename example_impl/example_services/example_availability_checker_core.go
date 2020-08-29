/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package example_services

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type ExampleAvailabilityCheckerCore struct{}

func (e ExampleAvailabilityCheckerCore) IsServiceUp(toCheck services.Service, dependencies []services.Service) bool {
	castedService := toCheck.(ExampleService)
	socket := castedService.GetHelloWorldSocket()
	url := fmt.Sprintf("http://%v:%v", socket.IPAddr, socket.Port)

	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}
	_, err := httpClient.Get(url)
	if err != nil {
		logrus.Tracef("Service not yet available due to the following error:")
		fmt.Fprintln(logrus.StandardLogger().Out, err)
		return false
	} else {
		return true
	}
}

func (e ExampleAvailabilityCheckerCore) GetTimeout() time.Duration {
	return 30 * time.Second
}

