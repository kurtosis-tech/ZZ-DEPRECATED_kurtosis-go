/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package testsuite_impl

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
)

type Testsuite struct {
	serviceImage string
}

func NewTestsuite(serviceImage string) *Testsuite {
	return &Testsuite{serviceImage: serviceImage}
}


func (suite Testsuite) GetTests() map[string]testsuite.Test {
	return map[string]testsuite.Test{
		"singleNodeNginxTest": SingleNodeExampleTest{ServiceImage: suite.serviceImage},
		"fixedSizeNginxTest": FixedSizeNginxTest{ServiceImage: suite.serviceImage},
	}
}

func (suite Testsuite) GetNetworkWidthBits() uint32 {
	return 8
}


