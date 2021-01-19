/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package impl

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
	"github.com/kurtosis-tech/kurtosis-go/lib_core/api/generated"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
)

type SuiteMetadataSerializingService struct {}

func (s SuiteMetadataSerializingService) SerializeSuiteMetadata(
		ctx context.Context,
		metadata *generated.TestSuiteMetadata) (*emptypb.Empty, error) {

	return nil, nil
}

// TODO Write tests for this by splitting it into metadata-generating function and writing function
//  then testing the metadata-generating
func (client KurtosisClient) printSuiteMetadataToFile(testSuite testsuite.TestSuite, filepath string) error {
	logrus.Debugf("Printing test suite metadata to file '%v'...", filepath)
	fp, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return stacktrace.Propagate(err, "No file exists at %v", filepath)
	}
	defer fp.Close()

	allTestMetadata := map[string]TestMetadata{}
	for testName, test := range testSuite.GetTests() {
		testConfig := test.GetTestConfiguration()

		// We create this "set" of used artifact URLs because the user could declare
		//  multiple artifacts with the same URL
		usedArtifactUrls := map[string]bool{}
		for _, artifactUrl := range testConfig.FilesArtifactUrls {
			usedArtifactUrls[artifactUrl] = true
		}

		artifactUrlsById := map[artifact_id_provider.ArtifactID]string{}
		for artifactUrl, _ := range usedArtifactUrls {
			artifactId, err := client.artifactIdProvider.GetArtifactId(artifactUrl)
			if err != nil {
				return stacktrace.Propagate(err, "An error occurred getting the artifact ID for URL '%v'", artifactUrl)
			}
			artifactUrlsById[artifactId] = artifactUrl
		}

		testMetadata := NewTestMetadata(
			testConfig.IsPartitioningEnabled,
			artifactUrlsById)
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
