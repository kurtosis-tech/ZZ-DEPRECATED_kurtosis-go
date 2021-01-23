/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/lib/execution"
	"github.com/kurtosis-tech/kurtosis-go/testsuite/execution_impl"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	successExitCode = 0
	failureExitCode = 1
)

func main() {
	customParamsJsonArg := flag.String(
		"custom-params-json",
		"{}",
		"JSON string containing custom data that the testsuite will deserialize to modify runtime behaviour",
	)

	kurtosisApiSocketArg := flag.String(
		"kurtosis-api-socket",
		"",
		"Socket in the form of address:port of the Kurtosis API container",
	)

	logLevelArg := flag.String(
		"log-level",
		"",
		"Loglevel string that the test suite will output with",
	)

	flag.Parse()

	// REPLACE WITH YOUR OWN CONFIGURATOR
	configurator := execution_impl.NewExampleTestsuiteConfigurator()

	suiteExecutor := execution.NewTestSuiteExecutor(*kurtosisApiSocketArg, *logLevelArg, *customParamsJsonArg, configurator)
	if err := suiteExecutor.Run(context.Background()); err != nil {
		logrus.Errorf("An error occurred running the test suite executor:")
		fmt.Fprintln(logrus.StandardLogger().Out, err)
		os.Exit(failureExitCode)
	}
	os.Exit(successExitCode)
}
