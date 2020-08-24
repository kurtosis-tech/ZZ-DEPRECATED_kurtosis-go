package services

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
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
	return err == nil
}

func (e ExampleAvailabilityCheckerCore) GetTimeout() time.Duration {
	return 30 * time.Second
}

