/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package client

// Package class, intended to get written to JSON, containing metadata about the test suite
type TestSuiteMetadata struct {
	TestNames map[string]bool

	NetworkWidthBits uint32
}
