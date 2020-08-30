/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package example_testsuite

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
)

// TODO example of parameterizing your image
type ExampleTestsuite struct {}

func (suite ExampleTestsuite) GetTests() map[string]testsuite.Test {
	return map[string]testsuite.Test{
		"singleNodeExampleTest": SingleNodeExampleTest{},
		"fixedSizeExampleTest": FixedSizeExampleTest{},
	}
}

func (suite ExampleTestsuite) GetNetworkWidthBits() uint32 {
	return 8
}


