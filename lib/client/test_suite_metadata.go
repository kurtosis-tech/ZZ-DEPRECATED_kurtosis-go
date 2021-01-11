/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package client

import "github.com/kurtosis-tech/kurtosis-go/lib/client/artifact_id_provider"

// ------------------------------------------------------------------------------------------------------
// NOTE: The fields of these classes need to be public for serialization to work, but we add constructors
//  anyways to try and reduce the oops-I-forgot-to-initialize-this-field errors that always exist because
//  Go frustratingly isn't object-oriented :|
// ------------------------------------------------------------------------------------------------------

type TestMetadata struct {
	IsPartitioningEnabled bool	`json:"isPartitioningEnabled"`

	// A map of all the artifacts that the test wants,
	//  which the initializer will download and make ready for the test at runtime
	// The map is in the form of ID -> URL, where the ID is:
	//	1. The ID that the initializer should associate the artifact with after downloading it and
	//	2. The ID that the client will use to retrieve the artifact when a test requests it
	UsedArtifacts map[string]string `json:"usedArtifacts"`
}

func NewTestMetadata(isPartitioningEnabled bool, artifactUrlsById map[artifact_id_provider.ArtifactID]string) *TestMetadata {
	var usedArtifacts = map[string]string{}
	for artifactId, url := range artifactUrlsById {
		usedArtifacts[string(artifactId)] = url
	}
	return &TestMetadata{
		IsPartitioningEnabled: isPartitioningEnabled,
		UsedArtifacts:         usedArtifacts,
	}
}


// Package class, intended to get written to JSON, containing metadata about the test suite
type TestSuiteMetadata struct {
	NetworkWidthBits uint32		`json:"networkWidthBits"`
	TestMetadata map[string]TestMetadata `json:"testMetadata"`
}

func NewTestSuiteMetadata(networkWidthBits uint32, testMetadata map[string]TestMetadata) *TestSuiteMetadata {
	return &TestSuiteMetadata{
		NetworkWidthBits: networkWidthBits,
		TestMetadata: testMetadata,
	}
}


