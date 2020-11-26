/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package main

import (
	"flag"
	"fmt"
	"github.com/kurtosis-tech/kurtosis-go/lib/client"
	"github.com/kurtosis-tech/kurtosis-go/testsuite/testsuite_impl"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	// ------------------- Kurtosis-internal params -------------------------------
	metadataFilepath := flag.String(
		"metadata-filepath",
		"",
		"The filepath of the file in which the test suite metadata should be written")
	testArg := flag.String(
		"test",
		"",
		"The name of the test to run")
	kurtosisApiIpArg := flag.String(
		"kurtosis-api-ip",
		"",
		"IP address of the Kurtosis API endpoint")
	logLevelArg := flag.String(
		"log-level",
		"",
		"String corresponding to Logrus log level that the test suite will output with",
		)
	servicesDirpathArg := flag.String(
		"services-relative-dirpath",
		"",
		"Dirpath, relative to the root of the suite execution volume, where directories for each service should be created")

	// -------------------- Testsuite-custom params ----------------------------------
	serviceImageArg := flag.String(
		"service-image",
		"",
		"Name of Nginx Docker image that will be used to launch service containers")


	flag.Parse()

	level, err := logrus.ParseLevel(*logLevelArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred parsing the log level string: %v\n", err)
		os.Exit(1)
	}
	logrus.SetLevel(level)

	testSuite := testsuite_impl.NewTestsuite(*serviceImageArg)
	exitCode := client.Run(testSuite, *metadataFilepath, *servicesDirpathArg, *testArg, *kurtosisApiIpArg)
	os.Exit(exitCode)
}
