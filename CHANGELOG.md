## 1.1.0
* Add Apache license
* Fix changelog check in CircleCI 
* Cleaning TODOs 
* Added a README pointing users to the main Kurtosis docs
* Cleaned up `build_and_run.sh` with lessons learned from upgrading the Avalanche test suite to Kurtosis 1.0
* Explain nil start command for the example impl
* Added a new bootstrapping process for creating Kurtosis Go testsuites from scratch
* Add [the comment-run Github Action](https://github.com/nwtgck/actions-comment-run/tree/20297f070391450752be7ac1ebd454fb53f62795#pr-merge-preview) to the repository in order to set up [a workaround for Github not passing secrets to untrusted PRs](https://github.community/t/secrets-for-prs-who-are-not-collaborators/17712), which would prevent auth'd Kurtosis from running
* Simplified the bootstrapping process quite a bit
* In addition to building `develop` and `master` images, build `X.Y.Z` tag images
* Cleaned up an over-aggressive check that was causing testsuite log-listing to fail
* When no arguments are provided to `build_and_run.sh`, the script errors
* In CircleCI config, don't run the `validate` workflow on `develop` and `master` (because they should already be validated by PR merge)

## 1.0.0
* Created example test suite to validate that the client library work
* Bugfix in volume-writing location, and force pretty formatting on written logs
* Made the existing test actually query the node it created
* Added another test to demonstrate an initial network setup
* Adding copyright headers
* Renamed tests to have more descriptive names
* When asked about test suite data, send back a JSON of test suite metadata (rather than just a list of test names)
* Made log level configurable
* Add CircleCI
* Upgraded example Go implementation to show the use of custom environment variables
* Build a Docker image on each merge to the develop branch
* Accept a new Docker parameter, `SERVICES_RELATIVE_DIRPATH`, for the location (relative to the suite execution volume root) where file IO for the services created during test execution
* Consolidate all the scripts into `build_and_run.sh` which will actually run the test suite for testing purposes
* Switch to `master` release track from Kurtosis core
