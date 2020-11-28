/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package services

import "time"

/*
The developer should implement their own use-case-specific interface that extends this one
 */
type Service interface {
	// Returns the IP address of the service
	GetIPAddress() string

	// Blocks until the service is available or the timeout is reached
	WaitForAvailability(timeout time.Duration) error
}

