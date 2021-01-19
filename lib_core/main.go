/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package main

import (
	"flag"
	"github.com/kurtosis-tech/kurtosis-go/lib_core/lib_core_process_consts"
)

func main() {
	portArg := flag.Int(
		lib_core_process_consts.PortFlag,
		0,
		"Port argument that the RPC server will be started on")

	paramsJsonArg := flag.String(
		lib_core_process_consts.ParamsJsonFlag,
		"",
		"JSON containing Kurtosis-specific parameters that control how the testsuite will execute")

	// Create Core object
	// TODO Launch RPC-listening server

	// TODO receive the JSON of the Kurtosis args
	// Port to listen on
	//
}
