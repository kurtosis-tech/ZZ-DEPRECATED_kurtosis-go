/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package services

/*
The developer should implement their own use-case-specific interface that extends this one
 */
type Service interface {
	// Returns the IP address of the service
	GetIPAddress() string

	// Returns true if the service is available
	IsAvailable() bool
}

