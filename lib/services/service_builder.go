/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package services

import "github.com/kurtosis-tech/kurtosis-go/lib/kurtosis_service"

type ServiceBuilder struct {
	// The core dictating initialization logic
	core *ServiceBuilderCore

	// The dirpath ON THE SUITE CONTAINER where the service-to-be's directory will be stored
	serviceDirpath string

	// The handle to manipulating the test environment
	kurtosisService *kurtosis_service.KurtosisService
}

func NewServiceBuilder(core *ServiceBuilderCore, serviceDirpath string, kurtosisService *kurtosis_service.KurtosisService) *ServiceBuilder {
	return &ServiceBuilder{core: core, serviceDirpath: serviceDirpath, kurtosisService: kurtosisService}
}



func (builder *ServiceBuilder) Build() {

}