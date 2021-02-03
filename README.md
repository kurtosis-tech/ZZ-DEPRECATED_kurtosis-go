DEPRECATION NOTICE
==================
This repo has been deprecated in favor of [Kurtosis Libs](https://github.com/kurtosis-tech/kurtosis-libs) as of 2021-02-02. This repo will stay around for a while to give time to migrate, but will eventually be archived. In order to migrate your testsuite repo using the old Kurtosis Go module across:

1. From the root of your testsuite repo, run `sed -i '' 's,github.com/kurtosis-tech/kurtosis-go,github.com/kurtosis-tech/kurtosis-libs/golang,g' $(find . -type f -name '*.go')` to swap the old Kurtosis Go module for the new Kurtosis Libs one in all your Go code
1. Run `go get github.com/kurtosis-tech/kurtosis-libs/golang@THE-VERSION-YOU-WERE-USING` (all version tags have been ported across to Kurtosis Libs, so you can use the same version)
1. Run `go mod tidy` to get rid of the old Kurtosis Go module import


Kurtosis Go Client
==================
This repo contains:

1. The Kurtosis Go client for using the Kurtosis platform with Golang
2. [An example implementation of the Go client](./testsuite), to demonstrate how to use the client
3. A generator for building new Kurtosis Go testsuites

### Quickstart
To get started building your own Kurtosis testsuite in Go, see [the Quickstart docs](https://docs.kurtosistech.com/quickstart.html).

### Library Development
Each library needs to talk with Kurtosis Core, and the Kurtosis Core API is defined via Protobuf. Rather than storing the Protobufs in Git submodules (which add significant complexity), the `.proto` files are simply copied from the relevant version of Kurtosis Core. In the future, we can move to a more productized solution.

To regenerate the bindings corresponding to the Protobuf files, use the `scripts/regenerate-protobuf-output.sh` script.
