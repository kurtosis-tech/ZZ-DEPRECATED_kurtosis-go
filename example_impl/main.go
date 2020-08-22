package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	testNamesFilepathArg := flag.String(
		"test-names-filepath",
		"",
		"The filepath of the file in which the names of all the tests should be written",
	)
	flag.Parse()

	// TODO replace with test suite

	fp, err := os.OpenFile(*testNamesFilepathArg, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Errorf("No file exists at %v", *testNamesFilepathArg)
		os.Exit(1)
	}
	defer fp.Close()

	testNames := []string{
		"test1",
		"test2",
		"test3",
	}
	for _, line := range testNames {
		fp.WriteString(line + "\n")
	}
}
