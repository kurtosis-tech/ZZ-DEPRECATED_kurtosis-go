/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package method_types

type RemoveServiceArgs struct {
	ServiceID string	`json:"serviceId"`
	ContainerStopTimeoutSeconds int `json:"containerStopTimeoutSeconds"`
}
