/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package execution

import "github.com/kurtosis-tech/kurtosis-go/lib/testsuite"

type TestSuiteConfigurator interface {
	SetLogLevel(logLevelStr string) error

	ParseParamsAndCreateSuite(paramsJsonStr string) (testsuite.TestSuite, error)
}
