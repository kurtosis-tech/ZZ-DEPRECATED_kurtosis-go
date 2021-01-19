/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package kurtosis

import (
	"github.com/palantir/stacktrace"
	"os"
	"os/exec"
	"strconv"
)

type executionAction int

const (
	printSuiteMetadataAction executionAction = iota
	runTestAction
)

func Execute(
		libCoreBinaryFilepath string,
		libCorePort int,
		kurtosisOptionsJson string,
		customOptionsJson string,
		suiteLogLevel string,
		customKurtosisCore KurtosisCore) error {
	// TODO Launch Kurtosis lib core as a separate coprocess, passing in the options JSON as-is
	libCoreLaunchingCmd := exec.Command(libCoreBinaryFilepath, strconv.Itoa(libCorePort))
	libCoreLaunchingCmd.Stdout = os.Stdout
	libCoreLaunchingCmd.Stderr = os.Stderr

	if err := libCoreLaunchingCmd.Start(); err != nil {
		return stacktrace.Propagate(err, "An error occurred starting the lib core coprocess")
	}

	// TODO Create RPC client for sending messages to the coprocess

	// TODO actually call GetExecutionAction
	executionAction := printSuiteMetadataAction

	if executionAction == printSuiteMetadataAction {

	}

	// TODO tell coprocess to kill itself

	if err := libCoreLaunchingCmd.Wait(); err != nil {
		return stacktrace.Propagate(err, "An error occurred waiting for the lib core coprocess to stop after being told to exit")
	}
}

