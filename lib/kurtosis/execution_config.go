/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package kurtosis

import "github.com/kurtosis-tech/kurtosis-go/lib/testsuite"

type ExecutionConfig interface {
	// Should take in the JSON of custom params intended for the testsuite, parse them, and create
	//  an instance of the user's testsuite
	CreateTestSuite(testSuiteParamsJson string) (*testsuite.TestSuite, error)
}
