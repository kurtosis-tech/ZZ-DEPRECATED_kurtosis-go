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
