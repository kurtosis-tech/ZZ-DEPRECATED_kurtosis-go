package client

import (
	"github.com/palantir/stacktrace"
	"strings"
)

func Run(testNamesFilepath string, test string) error {
	testNamesFilepath = strings.TrimSpace(testNamesFilepath)
	test = strings.TrimSpace(test)

	isTestNamesFilepathEmpty := len(testNamesFilepath) == 0
	isTestEmpty := len(test) == 0

	// Only one of these should be set; if both are set then it's an error
	if isTestNamesFilepathEmpty == isTestEmpty {
		return stacktrace.NewError("Exactly one of test-names-filepath and the test-name-to-run should be set")
	}

	/*
	if !isTestNamesFilepathEmpty {
		fp, err := os.OpenFile(testNamesFilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logrus.Errorf("No file exists at %v", testNamesFilepath)
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
	 */

	return nil
}