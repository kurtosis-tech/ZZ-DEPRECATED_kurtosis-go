package kurtosis_service

type AddServiceArgs struct {
	IPPlaceholder	      string 			`json:"ipPlaceholder"`
	ImageName             string            `json:"imageName"`
	UsedPorts             []string          `json:"usedPorts"`
	StartCmd              []string          `json:"startCommand"`
	DockerEnvironmentVars map[string]string `json:"dockerEnvironmentVars"`
	TestVolumeMountFilepath string			`json:"testVolumeMountFilepath"`
}

type AddServiceResponse struct {
	ServiceID string 	`json:"serviceId"`
	IPAddress string 	`json:"ipAddress"`
}

type RemoveServiceArgs struct {
	ServiceID string	`json:"serviceId"`
}

