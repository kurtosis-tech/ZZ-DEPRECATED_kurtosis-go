package client

// Package class, intended to get written to JSON, containing metadata about the test suite
type TestSuiteMetadata struct {
	TestNames map[string]bool

	NetworkWidthBits uint32
}
