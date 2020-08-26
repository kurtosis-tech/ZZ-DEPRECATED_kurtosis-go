package main

import (
	"flag"
	"github.com/kurtosis-tech/kurtosis-go/example_impl/testsuite"
	"github.com/kurtosis-tech/kurtosis-go/lib/client"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	testNamesFilepathArg := flag.String(
		"test-names-filepath",
		"",
		"The filepath of the file in which the names of all the tests should be written")
	testArg := flag.String(
		"test",
		"",
		"The name of the test to run")
	kurtosisApiIp := flag.String(
		"kurtosis-api-ip",
		"",
		"IP address of the Kurtosis API endpoint")
	flag.Parse()

	// TODO Make this parameterized
	logrus.SetLevel(logrus.TraceLevel)

	testSuite := testsuite.ExampleTestsuite{}

	exitCode := client.Run(testSuite, *testNamesFilepathArg, *testArg, *kurtosisApiIp)
	os.Exit(exitCode)
}
