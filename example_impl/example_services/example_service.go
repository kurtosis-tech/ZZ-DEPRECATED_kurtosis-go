/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package example_services

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
)

type ExampleService interface {
	services.Service

	GetHelloWorldSocket() Socket
}
