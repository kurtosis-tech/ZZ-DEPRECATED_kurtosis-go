/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package kurtosis

// Struct to contain information about the Docker container running the testsuite
//  that only the user knows
type SuiteContainerConfig interface {
	// Gets the path on the Docker image running the testsuite where the lib core binary
	//  lives
	GetLibCoreBinaryFilepath() string

	// Should take in the suite log level string, parse it, and set the appropriate log level for
	//  whichever logging framework is being used in the testsuite
	ConfigureLogLevel(suiteLogLevelStr string)
}
