/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package kurtosis

import (
	"context"
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/lib_core/api/generated"
	"github.com/kurtosis-tech/kurtosis-go/lib_core/lib_core_process_consts"
	"github.com/palantir/stacktrace"
	"google.golang.org/grpc"
	"os"
	"os/exec"
	"strconv"
)

type executionAction int

const (
	libCoreListenInterface = "127.0.0.1"

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
	ctx := context.Background()

	testsuite, err := executor.executionConfig.CreateTestSuite(suiteParamsJson)
	if err != nil {
		return stacktrace.Propagate(
			err,
			"An error occurred instantiating the testsuite using custom params JSON '%v'",
			suiteParamsJson)
	}

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

	libCoreListenStr := fmt.Sprintf("%v:%v", libCoreListenInterface, libCoreListenPort)
	conn, err := grpc.Dial(libCoreListenStr)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred dialing the lib core process on '%v'", libCoreListenStr)
	}
	defer conn.Close()

	wrapperExecutionPathClient := generated.NewWrapperExecutionPathServiceClient(conn)
	getPathResp, err := wrapperExecutionPathClient.GetPath(ctx, nil)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred making the call to lib core to get which execution path the wrapper lib should to go down")
	}
	path := getPathResp.GetPath()

	switch path {
	case generated.WrapperExecutionPath_PRINT_SUITE_METADATA:
		// TODO package up the test metadata
		// TODO ask lib core to print the suite metadata
	case generated.WrapperExecutionPath_RUN_TEST:
		// TODO create a NetworkContext wrapper (which really just forwards all its calls to lib core)
		// TODO get the test specified
	default:
		return stacktrace.NewError(
			"Lib core responded that the wrapper should go down unknown execution path '%v'; this is a code bug",
			path)
	}

	// TODO tell coprocess to kill itself

	if err := libCoreLaunchingCmd.Wait(); err != nil {
		return stacktrace.Propagate(err, "An error occurred waiting for the lib core coprocess to stop after being told to exit")
	}
}

