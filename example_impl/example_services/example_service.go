package example_services

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
)

type ExampleService interface {
	services.Service

	GetHelloWorldSocket() Socket
}
