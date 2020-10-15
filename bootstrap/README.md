My Kurtosis Testsuite
=====================
Welcome to your new Kurtosis testsuite skeleton! To get started, from the root directory:

1. `git init` to initialize this directory as a Git repo
1. `go get github.com/kurtosis-tech/kurtosis-go` to pull the Kurtosis Go client lib as a dependency
1. `git add .` to add all files
1. `git commit -m "Init commit"` to commit these files (don't skip this - it's used for running!)
1. `scripts/build_and_run.sh all` to build and run your test suite with the latest version of Kurtosis
1. Using the documentation on [the Kurtosis docs](https://github.com/kurtosis-tech/kurtosis-docs), modify the files within the `example_impl` to build a testsuite for your use case

Some helpful tips:
* If you rename the `example_impl` directory, you'll need to modify both the `Dockerfile` and the `build_and_run.sh` script
* `build_and_run.sh` can be run with the `build` argument to just build, `run` to just run, and `help` to see how to add Docker arguments
