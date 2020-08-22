package main

import (
	"flag"
	"github.com/kurtosis-tech/kurtosis-go/client"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	testNamesFilepathArg := flag.String(
		"test-names-filepath",
		"",
		"The filepath of the file in which the names of all the tests should be written",
	)
	testArg := flag.String(
		"test",
		"",
		"The name of the test to run",
	)

	flag.Parse()

	err := client.Run(*testNamesFilepathArg, *testArg)
	if err != nil {
		logrus.Errorf("An error occurred running the client: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
