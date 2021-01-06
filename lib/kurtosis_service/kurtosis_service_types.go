/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package kurtosis_service

type AddServiceArgs struct {
	ServiceID  		string 	`json:"serviceId"`
	PartitionID 	string 	`json:"partitionId"`
	IPPlaceholder	string 	`json:"ipPlaceholder"`
	ImageName       string  `json:"imageName"`

	// This is in Docker port specification syntax, e.g. "80" (default TCP) or "80/udp"
	// It might even support ranges (e.g. "90:100/tcp"), though this is untested as of 2020-12-08
	UsedPorts             []string          `json:"usedPorts"`

	StartCmd              []string          `json:"startCommand"`
	DockerEnvironmentVars map[string]string `json:"dockerEnvironmentVars"`
	TestVolumeMountDirpath string			`json:"testVolumeMountDirpath"`
}

type AddServiceResponse struct {
	IPAddress string 	`json:"ipAddress"`
}

type RemoveServiceArgs struct {
	ServiceID string	`json:"serviceId"`
	ContainerStopTimeoutSeconds int `json:"containerStopTimeoutSeconds"`
}

type RegisterTestExecutionArgs struct {
	TestTimeoutSeconds int
}

