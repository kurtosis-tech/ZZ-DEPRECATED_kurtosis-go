/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package kurtosis

type KurtosisCore interface {
	ConfigureLogLevel(logLevelStr string)

	CreateTestSuite(customOptionsJson string)
}
