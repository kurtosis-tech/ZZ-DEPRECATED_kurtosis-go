package client

import (
	"github.com/kurtosis-tech/kurtosis-go/kurtosis_service"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

func Run(testNamesFilepath string, test string, kurtosisApiIp string) error {
	testNamesFilepath = strings.TrimSpace(testNamesFilepath)
	test = strings.TrimSpace(test)

	isTestNamesFilepathEmpty := len(testNamesFilepath) == 0
	isTestEmpty := len(test) == 0

	// Only one of these should be set; if both are set then it's an error
	if isTestNamesFilepathEmpty == isTestEmpty {
		return stacktrace.NewError("Exactly one of test-names-filepath and the test-name-to-run should be set")
	}

	if !isTestNamesFilepathEmpty {
		if err := printTestsToFile(testNamesFilepath); err != nil {
			return stacktrace.Propagate(err, "An error occurred printing tests to file '%v'", testNamesFilepath)
		}
	} else if !isTestEmpty {
		// TODO parameterize
		if err := runTest(test, kurtosisApiIp); err != nil {
			return stacktrace.Propagate(err, "An error occurred running test '%v'", test)
		}
	}

	return nil
}

// =========================================== Private helper functions ========================================
func printTestsToFile(testNamesFilepath string) error {
	logrus.Debugf("Printing tests to file '%v'...", testNamesFilepath)
	fp, err := os.OpenFile(testNamesFilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return stacktrace.Propagate(err, "No file exists at %v", testNamesFilepath)
	}
	defer fp.Close()

	testNames := []string{
		"test1",
		"test2",
		"test3",
	}
	for _, line := range testNames {
		fp.WriteString(line + "\n")
	}

	return nil
}

func runTest(test string, kurtosisApiIp string) error {
	kurtosisService := kurtosis_service.NewKurtosisService(kurtosisApiIp)

	// TODO replace with parameterized test execution
	if err := kurtosisService.RegisterTestExecution(60); err != nil {
		return stacktrace.Propagate(err, "An error occurred registering the test execution with the API container")
	}

	_, containerId, err := kurtosisService.AddService(
		"nginxdemos/hello",
		map[int]bool{
			80: true,
		},
		"BLAHBLAH",
		[]string{},
		map[string]string{},
		"/nothing-yet")
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred adding a new service")
	}

	time.Sleep(10 * time.Second)

	if err := kurtosisService.RemoveService(containerId, 30); err != nil {
		return stacktrace.Propagate(err, "An error occurred removing the new service")
	}

	return nil
}