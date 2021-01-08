/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package client

type TestMetadata struct {
	IsPartitioningEnabled bool	`json:"isPartitioningEnabled"`
}

// Package class, intended to get written to JSON, containing metadata about the test suite
type TestSuiteMetadata struct {
	NetworkWidthBits uint32		`json:"networkWidthBits"`

	TestMetadata map[string]TestMetadata	`json:"testMetadata"`
}
