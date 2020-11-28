## TBD
* Remove socket in favor of `ExampleService.GetIpAddress` and `ExapleService.GetPort` methods
* Remove TODO on allowing non-TCP ports
* Removed the `example_` prefix to make bootstrapping even easier (users need only modify the existing suite, no need to remove the `example_` prefix)
* Heavily refactored the client architecture to make it much less confusing to define testsuite infrastructure:
    * **Detailed changes:**
        * The notion of `dependencies` that showed up in several places (e.g. `ServiceInitializerCore.GetStartCommand`, `ServiceAvailabilityCheckerCore.IsServiceUp`, etc) have been removed due to being too confusing
        * Services: 
            * The `Service` interface has received two new methods, `GetIPAddress` and `IsAvailable` to better reflect what services are
            * `ServiceInitializerCore`, `ServiceInitializer`, and `ServiceAvailabilityCheckerCore` have been removed
            * `ServiceInitializerCore`'s functionality has been subsumed by a new interface, `DockerContainerInitializer`
            * `ServiceAvailabilityChecker` renamed to `AvailabilityChecker`
            * The old `ServiceAvailabilityChecker.WaitForStartup` method is now `AvailabilityChecker.WaitForStartup(timeBetweenPolls time.Duration, maxNumRetries int)`
        * Networks: 
            * `ServiceNetwork` has been renamed `NetworkContext`, with `NetworkContext.AddService(DockerContainerInitializer) (Service, AvailabilityChecker, error)` replacing the old `ServiceNetwork.AddService(ConfigurationID, ServiceID, map[ServiceID]bool) (*ServiceAvailabilityChecker, error)` method
            * Test networks are no longer instantiated in two separate configuration/instantiation phases, and are simply instantiated with a `Test.Setup` method
            * The notion of "service configuration" that was used during the network configuration phase has been removed, since networks are simply instantiated now
            * `ServiceNetworkBuilder` has been removed
            * `NetworkLoader` has been removed
        * Testsuite:
            * The `Test.GetNetworkLoader` method has been replaced with `Test.Setup(NetworkContext) Network`
                * The `Network` return type is still `interface{}`, so users can return `NetworkContext` directly or wrap it in a more test-friendly custom object
    * **Remediation instructions:**
        1. Services:
            1. Update any existing implementations of the `Service` interface to implement the `GetIPAddress` and `IsAvailable` methods
                * `GetIPAddress` should return an IP address `string` that's given to the `Service` at construction time
                * `IsAvailable` should contain the logic that's currently inside the `ServiceAvailabilityCheckerCore` implementation for the service
            1. Delete any `ServiceAvailabilityCheckerCore` implementations after their logic has been moved to `Service.IsAvailable`
            1. For each existing `ServiceInitializerCore` implementation, modify it to instead implement `DockerContainerInitializer`:
                * Rename the file/struct if desired
                * Add a `GetDockerImage` method (which will likely just return a `string` passed in at construction time)
                * Remove the `dependencies` parameter from the `InitializeMountedFiles` function
                * Remove the `dependencies` parameter from the `GetStartCommand` function
        1. Networks:
            1. For each `Network` implementation, 
            1. Modify all uses of the old `ServiceAvailabilityCheckerCore.WaitForStartup()` method to match the new `AvailabilityChecker.WaitForStartup(timeBetweenPolls time.Duration, maxNumRetries int)` signature
        1. Testsuite:
            1. For each old `Test` implementation, update it to use the updated `Test` interface:
                1. Add a `Setup` method to match the new interface
                1. Move the logic from the test's network loader's `NetworkLoader.Initialize` function to the test's `Setup` method (or extract it into a new struct if you prefer)
                1. Remove the `GetNetworkLoader` method
        * Networks:
            1. Delete each `NetworkLoader` implementation after verifying that the `Initializer` contents are correctly being called 
        * All methods on `ServiceInitializerCore` should be 
        * ServiceInitializerCore
        * TODO instructions

## 1.1.1
* Remove log filepath (which is no longer needed now that Kurtosis core reads Docker logs directly)
* Switch to using [our forked version of action-comment-run](https://github.com/mieubrisse/actions-comment-run) that allows user whitelisting
* Bump kurtosis-core to 1.1.0
* Make the requests to the Kurtosis API container retry every second, with 10s retry maximum for normal operations (e.g. add/remove services) and 60s retry maximum for test suite registration
* Update the version of the `actions-comment-run` Github Action which allows for running CI on untrusted PRs, to match the advice we give in the "Running In CI" instructions

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
