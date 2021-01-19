/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package kurtosis

import (
	"github.com/kurtosis-tech/kurtosis-go/lib_core/lib_core_process_consts"
	"github.com/palantir/stacktrace"
	"os"
	"os/exec"
	"strconv"
)

type executionAction int

const (
	printSuiteMetadataAction executionAction = iota
	runTestAction

	// TODO can make this configurable if it causes problems for the user's testsuite
	libCoreListenPort = 5822
)

type KurtosisExecutor struct {
	suiteContainerConfig SuiteContainerConfig
	executionConfig ExecutionConfig
}



// Runs Kurtosis
// Intended to be a simple call that just consumes the arguments passed in to main
func (executor KurtosisExecutor) Execute(
		libCoreParamsJson string,
		suiteParamsJson string,
		suiteLogLevel string) error {
	// Launch the core lib as a separate process, communicable via gRPC
	libCoreBinaryFilepath := executor.suiteContainerConfig.GetLibCoreBinaryFilepath()
	libCoreLaunchingCmd := exec.Command(
		libCoreBinaryFilepath,
		"-" + lib_core_process_consts.PortFlag,
		strconv.Itoa(libCoreListenPort),
		"-" + lib_core_process_consts.ParamsJsonFlag,
		libCoreParamsJson)
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

