package example_testsuite

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
)

// TODO example of parameterizing your image
type ExampleTestsuite struct {}

func (suite ExampleTestsuite) GetTests() map[string]testsuite.Test {
	return map[string]testsuite.Test{
		"exampleTest1": ExampleTest1{},
		"exampleTest2": ExampleTest2{},
	}
}

func (suite ExampleTestsuite) GetNetworkWidthBits() uint32 {
	return 8
}


