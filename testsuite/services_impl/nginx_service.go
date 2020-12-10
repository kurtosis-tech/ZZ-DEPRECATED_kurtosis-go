/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package services_impl

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	nginxServicePort = 80
)

type NginxService struct{
	IPAddr string
}

func (e NginxService) GetIPAddress() string {
	return e.IPAddr
}

func (e NginxService) GetPort() int {
	return nginxServicePort
}

func (e NginxService) IsAvailable() bool {
	url := fmt.Sprintf("http://%v:%v", e.GetIPAddress(), nginxServicePort)
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		logrus.Tracef("Service not yet available due to the following error:")
		fmt.Fprintln(logrus.StandardLogger().Out, err)
		return false
	}
	return resp.StatusCode == 200
}

