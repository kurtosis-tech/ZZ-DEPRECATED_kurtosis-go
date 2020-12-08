/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package services_impl

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
)

type NginxService interface {
	services.Service

	GetIpAddress() string

	GetPort() int
}
