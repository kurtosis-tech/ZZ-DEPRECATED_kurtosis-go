/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package example_testsuite

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
)

type ExampleTestsuite struct {
	serviceImage string
}

func NewExampleTestsuite(serviceImage string) *ExampleTestsuite {
	return &ExampleTestsuite{serviceImage: serviceImage}
}


func (suite ExampleTestsuite) GetTests() map[string]testsuite.Test {
	return map[string]testsuite.Test{
		"singleNodeExampleTest": SingleNodeExampleTest{ServiceImage: suite.serviceImage},
		"fixedSizeExampleTest": FixedSizeExampleTest{ServiceImage: suite.serviceImage},
	}
}

func (suite ExampleTestsuite) GetNetworkWidthBits() uint32 {
	return 8
}


