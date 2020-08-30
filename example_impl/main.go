/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package main

import (
	"flag"
	"github.com/kurtosis-tech/kurtosis-go/example_impl/example_testsuite"
	"github.com/kurtosis-tech/kurtosis-go/lib/client"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	metadataFilepath := flag.String(
		"metadata-filepath",
		"",
		"The filepath of the file in which the test suite metadata should be written")
	testArg := flag.String(
		"test",
		"",
		"The name of the test to run")
	kurtosisApiIp := flag.String(
		"kurtosis-api-ip",
		"",
		"IP address of the Kurtosis API endpoint")
	flag.Parse()

	// TODO Make this parameterized
	logrus.SetLevel(logrus.TraceLevel)

	testSuite := example_testsuite.ExampleTestsuite{}

	exitCode := client.Run(testSuite, *metadataFilepath, *testArg, *kurtosisApiIp)
	os.Exit(exitCode)
}
