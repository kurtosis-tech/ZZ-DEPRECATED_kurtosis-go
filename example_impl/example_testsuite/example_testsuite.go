package example_testsuite

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
)

// TODO example of parameterizing your image
type ExampleTestsuite struct {}

func (e ExampleTestsuite) GetTests() map[string]testsuite.Test {
	return map[string]testsuite.Test{
		"singleNodeExampleTest": SingleNodeExampleTest{},
		"fixedSizeExampleTest": FixedSizeExampleTest{},
	}
}

