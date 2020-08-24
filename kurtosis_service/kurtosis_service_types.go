package kurtosis_service

type AddServiceArgs struct {
	IPPlaceholder	      string 			`json:"ipPlaceholder"`
	ImageName             string            `json:"imageName"`
	UsedPorts             []int          `json:"usedPorts"`
	StartCmd              []string          `json:"startCommand"`
	DockerEnvironmentVars map[string]string `json:"dockerEnvironmentVars"`
	TestVolumeMountFilepath string			`json:"testVolumeMountFilepath"`
}

type AddServiceResponse struct {
	ContainerID string 	`json:"containerId"`
	IPAddress string 	`json:"ipAddress"`
}

type RemoveServiceArgs struct {
	ContainerID string	`json:"containerId"`
	ContainerStopTimeoutSeconds int `json:"containerStopTimeoutSeconds"`
}

type RegisterTestExecutionArgs struct {
	TestTimeoutSeconds int
}

