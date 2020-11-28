/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package services_impl

const (
	nginxServicePort = 80
)
type NginxServiceImpl struct{
	IPAddr string
}

func (e NginxServiceImpl) GetIpAddress() string {
	return e.IPAddr
}

func (e NginxServiceImpl) GetPort() int {
	return nginxServicePort
}
