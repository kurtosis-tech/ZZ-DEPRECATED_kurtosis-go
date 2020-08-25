package client

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/lib/kurtosis_service"
	"github.com/kurtosis-tech/kurtosis-go/lib/networks"
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

const (
	errorExitCode = 1
	successExitCode = 0

	// TODO parameterize this!!
	testVolumeMountLocation = "/test-volume"
)

func Run(testSuite testsuite.TestSuite, testNamesFilepath string, testName string, kurtosisApiIp string) int {
	testNamesFilepath = strings.TrimSpace(testNamesFilepath)
	testName = strings.TrimSpace(testName)

	isTestNamesFilepathEmpty := len(testNamesFilepath) == 0
	isTestEmpty := len(testName) == 0

	// Only one of these should be set; if both are set then it's an error
	if isTestNamesFilepathEmpty == isTestEmpty {
		logrus.Error("Exactly one of test-names-filepath and the test-name-to-run should be set")
		return errorExitCode
	}

	if !isTestNamesFilepathEmpty {
		if err := printTestsToFile(testSuite, testNamesFilepath); err != nil {
			logrus.Errorf("An error occurred printing tests to file '%v':", testNamesFilepath)
			fmt.Fprintln(logrus.StandardLogger().Out, err)
			return errorExitCode
		}
	} else if !isTestEmpty {
		if err := runTest(testSuite, testName, kurtosisApiIp); err != nil {
			logrus.Errorf("An error occurred running test '%v':", testName)
			fmt.Fprintln(logrus.StandardLogger().Out, err)
			return errorExitCode
		}
	}
	return successExitCode
}

// =========================================== Private helper functions ========================================
func printTestsToFile(testSuite testsuite.TestSuite, testNamesFilepath string) error {
	logrus.Debugf("Printing tests to file '%v'...", testNamesFilepath)
	fp, err := os.OpenFile(testNamesFilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return stacktrace.Propagate(err, "No file exists at %v", testNamesFilepath)
	}
	defer fp.Close()

	for testName, _ := range testSuite.GetTests() {
		fp.WriteString(testName + "\n")
	}

	return nil
}

/*
Runs the single given test from the testsuite

Returns:
	setupErr: Indicates an error setting up the test that prevented the test from running
	testErr: Indicates an error in the test itself, indicating a test failure
*/
func runTest(testSuite testsuite.TestSuite, testName string, kurtosisApiIp string) error {
	kurtosisService := kurtosis_service.NewKurtosisService(kurtosisApiIp)

	tests := testSuite.GetTests()
	test, found := tests[testName]
	if !found {
		return stacktrace.NewError("No test in the test suite named '%v'", testName)
	}

	// Kick off a timer with the API in case there's an infinite loop in the user code that causes the test to hang forever
	hardTestTimeout := test.GetExecutionTimeout() + test.GetSetupBuffer()
	hardTestTimeoutSeconds := int(hardTestTimeout.Seconds())
	if err := kurtosisService.RegisterTestExecution(hardTestTimeoutSeconds); err != nil {
		return stacktrace.Propagate(err, "An error occurred registering the test execution with the API container")
	}

	logrus.Info("Configuring test network...")
	builder := networks.NewServiceNetworkBuilder(
		kurtosisService,
		testVolumeMountLocation)
	networkLoader, err := test.GetNetworkLoader()
	if err != nil {
		return stacktrace.Propagate(err, "Could not get network loader")
	}
	if err := networkLoader.ConfigureNetwork(builder); err != nil {
		return stacktrace.Propagate(err, "Could not configure test network")
	}
	network := builder.Build()
	logrus.Info("Test network configured")

	logrus.Info("Initializing test network...")
	availabilityCheckers, err := networkLoader.InitializeNetwork(network);
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred initialized the network to its starting state")
	}
	logrus.Info("Test network initialized")

	// Second pass: wait for all services to come up
	logrus.Info("Waiting for test network to become available...")
	for serviceId, availabilityChecker := range availabilityCheckers {
		logrus.Debugf("Waiting for service %v to become available...", serviceId)
		if err := availabilityChecker.WaitForStartup(); err != nil {
			return stacktrace.Propagate(err, "An error occurred waiting for service with ID %v to start up", serviceId)
		}
		logrus.Debugf("Service %v is available", serviceId)
	}
	logrus.Info("Test network is available")

	logrus.Info("Wrapping untyped network in user-custom type...")
	untypedNetwork, err := networkLoader.WrapNetwork(network)
	if err != nil {
		return stacktrace.Propagate(err, "Error occurred wrapping network in user-defined network type")
	}
	logrus.Info("Untyped network wrapped in user-custom type")

	logrus.Infof("Executing test '%v'...", testName)
	testResultChan := make(chan error)

	go func() {
		testResultChan <- runTestInGoroutine(test, untypedNetwork)
	}()

	// Time out the test so a poorly-written test doesn't run forever
	testTimeout := test.GetExecutionTimeout()
	var timedOut bool
	var testResultErr error
	select {
	case testResultErr = <- testResultChan:
		logrus.Tracef("Test returned result before timeout: %v", testResultErr)
		timedOut = false
	case <- time.After(testTimeout):
		logrus.Tracef("Hit timeout %v before getting a result from the test", testTimeout)
		timedOut = true
	}
	logrus.Tracef("After running test w/timeout: resultErr: %v, timedOut: %v", testResultErr, timedOut)

	if timedOut {
		return stacktrace.NewError("Timed out after %v waiting for test to complete", testTimeout)
	}
	logrus.Infof("Executed test '%v'", testName)

	if testResultErr != nil {
		return stacktrace.Propagate(testResultErr, "An error occurred when running the test")
	}

	return nil
}

// Little helper function meant to be run inside a goroutine that runs the test
func runTestInGoroutine(test testsuite.Test, untypedNetwork interface{}) (resultErr error) {
	// See https://medium.com/@hussachai/error-handling-in-go-a-quick-opinionated-guide-9199dd7c7f76 for details
	defer func() {
		if recoverResult := recover(); recoverResult != nil {
			logrus.Tracef("Caught panic while running test: %v", recoverResult)
			resultErr = recoverResult.(error)
		}
	}()
	test.Run(untypedNetwork, testsuite.TestContext{})
	logrus.Tracef("Test completed successfully")
	return
}
