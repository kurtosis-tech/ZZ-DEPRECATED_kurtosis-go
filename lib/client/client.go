/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package client

import (
	"crypto"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/lib/kurtosis_service"
	"github.com/kurtosis-tech/kurtosis-go/lib/networks"
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"

	// This is a special type of import that includes the correct hashing algorithm that we use
	// If we don't have the "_" in front, Goland will complain it's unused
	_ "golang.org/x/crypto/sha3"
)

const (
	errorExitCode = 1
	successExitCode = 0

	// NOTE: right now this is hardcoded in the initializer as part of the contract between Kurtosis & a test suite image -
	//  all test suite images MUST have this path available for mounting
	suiteExecutionVolumeMountDirpath = "/suite-execution"

	hashFunction = crypto.SHA3_256
)

func Run(testSuite testsuite.TestSuite, metadataFilepath string, servicesRelativeDirpath string, testName string, kurtosisApiIp string) int {
	// Only one of these should be set; if both are set then it's an error
	metadataFilepath = strings.TrimSpace(metadataFilepath)
	testName = strings.TrimSpace(testName)
	isMetadataFilepathEmpty := len(metadataFilepath) == 0
	isTestEmpty := len(testName) == 0
	if isMetadataFilepathEmpty == isTestEmpty {
		logrus.Error("Exactly one of 'metadata filepath' or 'test name to run' should be set")
		return errorExitCode
	}

	if !isMetadataFilepathEmpty {
		if err := printSuiteMetadataToFile(testSuite, metadataFilepath); err != nil {
			logrus.Errorf("An error occurred writing test suite metadata to file '%v':", metadataFilepath)
			fmt.Fprintln(logrus.StandardLogger().Out, err)
			return errorExitCode
		}
	} else if !isTestEmpty {
		servicesRelativeDirpath = strings.TrimSpace(servicesRelativeDirpath)
		if len(servicesRelativeDirpath) == 0 {
			logrus.Error("Services relative dirpath argument was empty")
			return errorExitCode
		}
		kurtosisApiIp = strings.TrimSpace(kurtosisApiIp)
		if len(kurtosisApiIp) == 0 {
			logrus.Error("Kurtosis API container IP argument was empty")
			return errorExitCode
		}
		if err := runTest(servicesRelativeDirpath, testSuite, testName, kurtosisApiIp); err != nil {
			logrus.Errorf("An error occurred running test '%v':", testName)
			fmt.Fprintln(logrus.StandardLogger().Out, err)
			return errorExitCode
		}
	}
	return successExitCode
}

// =========================================== Private helper functions ========================================
// TODO Write tests for this by splitting it into metadata-generating function and writing function
//  then testing the metadata-generating
func printSuiteMetadataToFile(testSuite testsuite.TestSuite, filepath string) error {
	logrus.Debugf("Printing test suite metadata to file '%v'...", filepath)
	fp, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return stacktrace.Propagate(err, "No file exists at %v", filepath)
	}
	defer fp.Close()

	allTestMetadata := map[string]TestMetadata{}
	for testName, test := range testSuite.GetTests() {
		testConfig := test.GetTestConfiguration()

		// "Set" of used artifact URLs
		usedArtifactUrls := map[string]bool{}
		for _, artifactUrl := range testConfig.FilesArtifactUrls {
			usedArtifactUrls[artifactUrl] = true
		}

		artifactUrlsByHash := map[string]string{}
		for artifactUrl, _ := range usedArtifactUrls {
			hexEncodedHash, err := hashArtifactUrl(artifactUrl)
			if err != nil {
				return stacktrace.Propagate(err, "An error occurred hashing artifact URL '%v'", artifactUrl)
			}
			artifactUrlsByHash[hexEncodedHash] = artifactUrl
		}

		testMetadata := NewTestMetadata(
			testConfig.IsPartitioningEnabled,
			artifactUrlsByHash)
		allTestMetadata[testName] = *testMetadata
	}
	suiteMetadata := TestSuiteMetadata{
		TestMetadata:     allTestMetadata,
		NetworkWidthBits: testSuite.GetNetworkWidthBits(),
	}

	bytes, err := json.Marshal(suiteMetadata)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred serializing test suite metadata to JSON")
	}

	if _, err := fp.Write(bytes); err != nil {
		return stacktrace.Propagate(err, "An error occurred writing the JSON string to file")
	}

	return nil
}

func hashArtifactUrl(artifactUrl string) (hexStr string, resultErr error) {
	hasher := hashFunction.New()
	artifactUrlBytes := []byte(artifactUrl)
	if _, err := hasher.Write(artifactUrlBytes); err != nil {
		return "", stacktrace.Propagate(err, "An error occurred writing the artifact URL to the hash function")
	}
	hexEncodedHash := hex.EncodeToString(hasher.Sum(nil))
	return hexEncodedHash, nil
}

/*
Runs the single given test from the testsuite

Args:
	servicesRelativeDirpath: Dirpath where per-service directories live, relative to the root of the suite execution volume
	testSuite: Test suite to run
	testName: Name of test to run
	kurtosisApiIp: IP address of the Kurtosis API container

Returns:
	setupErr: Indicates an error setting up the test that prevented the test from running
	testErr: Indicates an error in the test itself, indicating a test failure
*/
func runTest(servicesRelativeDirpath string, testSuite testsuite.TestSuite, testName string, kurtosisApiIp string) error {
	kurtosisService := kurtosis_service.NewDefaultKurtosisService(kurtosisApiIp)

	tests := testSuite.GetTests()
	test, found := tests[testName]
	if !found {
		return stacktrace.NewError("No test in the test suite named '%v'", testName)
	}

	// Kick off a timer with the API in case there's an infinite loop in the user code that causes the test to hang forever
	hardTestTimeout := test.GetExecutionTimeout() + test.GetSetupTeardownBuffer()
	hardTestTimeoutSeconds := int(hardTestTimeout.Seconds())
	if err := kurtosisService.RegisterTestExecution(hardTestTimeoutSeconds); err != nil {
		return stacktrace.Propagate(err, "An error occurred registering the test execution with the API container")
	}

	networkCtx := networks.NewNetworkContext(kurtosisService, suiteExecutionVolumeMountDirpath, servicesRelativeDirpath)

	logrus.Info("Setting up the test network...")
	untypedNetwork, err := test.Setup(networkCtx)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred setting up the test network")
	}
	logrus.Info("Test network set up")

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
