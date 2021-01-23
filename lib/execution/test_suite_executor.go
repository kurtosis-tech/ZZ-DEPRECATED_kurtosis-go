/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package execution

import (
	"context"
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/lib/core_api/bindings"
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

const (
	maxSuiteRegistrationRetries = 20
	timeBetweenSuiteRegistrationRetries = 500 * time.Millisecond
)

type TestSuiteExecutor struct {
	kurtosisApiSocket string
	logLevelStr string
	paramsJsonStr string
	configurator TestSuiteConfigurator
}

func NewTestSuiteExecutor(kurtosisApiSocket string, logLevelStr string, paramsJsonStr string, configurator TestSuiteConfigurator) *TestSuiteExecutor {
	return &TestSuiteExecutor{kurtosisApiSocket: kurtosisApiSocket, logLevelStr: logLevelStr, paramsJsonStr: paramsJsonStr, configurator: configurator}
}

func (executor *TestSuiteExecutor) Run(ctx context.Context) error {
	if err := executor.configurator.SetLogLevel(executor.logLevelStr); err != nil {
		return stacktrace.Propagate(err, "An error occurred setting the loglevel before running the testsuite")
	}

	suite, err := executor.configurator.ParseParamsAndCreateSuite(executor.paramsJsonStr)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred parsing the suite params JSON and creating the testsuite")
	}

	// TODO SECURITY: Use HTTPS to ensure you're conecting with real Kurtosis API servers
	conn, err := grpc.Dial(executor.kurtosisApiSocket, grpc.WithInsecure())
	if err != nil {
		return stacktrace.Propagate(
			err,
			"An error occurred creating a connection to the Kurtosis API server at '%v'",
			executor.kurtosisApiSocket)
	}
	defer conn.Close()

	suiteRegistrationClient := bindings.NewSuiteRegistrationServiceClient(conn)

	var suiteRegistrationResp *bindings.SuiteRegistrationResponse
	suiteRegistrationAttempts := 0
	for {
		if suiteRegistrationAttempts >= maxSuiteRegistrationRetries {
			return stacktrace.NewError(
				"Failed to register testsuite with API container, even after %v retries spaced %v apart",
				maxSuiteRegistrationRetries,
				timeBetweenSuiteRegistrationRetries)
		}

		resp, err := suiteRegistrationClient.RegisterSuite(ctx, &emptypb.Empty{})
		if err == nil {
			suiteRegistrationResp = resp
			break
		}
		logrus.Debugf("The following error occurred registering testsuite with API container; retrying in %v:", timeBetweenSuiteRegistrationRetries)
		fmt.Fprintln(logrus.StandardLogger().Out, err)
		time.Sleep(timeBetweenSuiteRegistrationRetries)
		suiteRegistrationAttempts++
	}

	action := suiteRegistrationResp.SuiteAction
	switch action {
	case bindings.SuiteAction_SERIALIZE_SUITE_METADATA:
		if err := runSerializeSuiteMetadataFlow(ctx, suite, conn); err != nil {
			return stacktrace.Propagate(err, "An error occurred running the suite metadata serialization flow")
		}
		return nil
	case bindings.SuiteAction_EXECUTE_TEST:
		// TODO run nserialize suite metadata flow
		return stacktrace.NewError("NOT IMPLEMENTED YET")
	default:
		return stacktrace.NewError("Encountered unrecognized action '%v'; this is a code bug", action)
	}
}

func runSerializeSuiteMetadataFlow(ctx context.Context, testsuite testsuite.TestSuite, conn *grpc.ClientConn) error {
	allTestMetadata := map[string]*bindings.TestMetadata{}
	for testName, test := range testsuite.GetTests() {
		testConfig := test.GetTestConfiguration()
		usedArtifactUrls := map[string]bool{}
		for _, artifactUrl := range testConfig.FilesArtifactUrls {
			usedArtifactUrls[artifactUrl] = true
		}

		testMetadata := &bindings.TestMetadata{
			IsPartitioningEnabled: testConfig.IsPartitioningEnabled,
			UsedArtifactUrls:      usedArtifactUrls,
		}
		allTestMetadata[testName] = testMetadata
	}

	networkWidthBits := testsuite.GetNetworkWidthBits()
	testSuiteMetadata := &bindings.TestSuiteMetadata{
		TestMetadata:     allTestMetadata,
		NetworkWidthBits: networkWidthBits,
	}

	metadataSerializationClient := bindings.NewSuiteMetadataSerializationServiceClient(conn)
	if _, err := metadataSerializationClient.SerializeSuiteMetadata(ctx, testSuiteMetadata); err != nil {
		return stacktrace.Propagate(err, "An error occurred sending the suite metadata to the Kurtosis API server")
	}

	return nil
}
