/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package network_partition_test

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/networks"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
	"github.com/kurtosis-tech/kurtosis-go/testsuite/services_impl/api"
	"github.com/kurtosis-tech/kurtosis-go/testsuite/services_impl/datastore"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	apiPartitionId networks.PartitionID = "api"
	datastorePartitionId networks.PartitionID = "datastore"

	datastoreServiceId services.ServiceID = "datastore"
	apiServiceId services.ServiceID = "api"

	waitForStartupTimeBetweenPolls = 1 * time.Second
	waitForStartupMaxNumPolls = 30

	testPersonId = 46
)

type NetworkPartitionTest struct {
	datstoreImage string
	apiImage string
}

func NewNetworkPartitionTest(datstoreImage string, apiImage string) *NetworkPartitionTest {
	return &NetworkPartitionTest{datstoreImage: datstoreImage, apiImage: apiImage}
}

// Instantiates the network with no partition and one person in the datatstore
func (test NetworkPartitionTest) Setup(networkCtx *networks.NetworkContext) (networks.Network, error) {
	datastoreInitializer := datastore.NewDatastoreContainerInitializer(test.datstoreImage)
	uncastedDatastoreSvc, datastoreChecker, err := networkCtx.AddService(datastoreServiceId, datastoreInitializer)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred adding the datastore service")
	}
	if err := datastoreChecker.WaitForStartup(waitForStartupTimeBetweenPolls, waitForStartupMaxNumPolls); err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred waiting for the datastore service to start")
	}

	// Go doesn't have generics so we need to do this cast
	datastoreSvc := uncastedDatastoreSvc.(*datastore.DatastoreService)

	apiInitializer := api.NewApiContainerInitializer(test.apiImage, datastoreSvc)
	uncastedApiSvc, apiChecker, err := networkCtx.AddService(apiServiceId, apiInitializer)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred adding the API service")
	}
	if err := apiChecker.WaitForStartup(waitForStartupTimeBetweenPolls, waitForStartupMaxNumPolls); err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred waiting for the API service to start")
	}

	apiSvc := uncastedApiSvc.(*api.ApiService)
	if err := apiSvc.AddPerson(testPersonId); err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred adding the test person in preparation for the test")
	}
	if err := apiSvc.IncrementBooksRead(testPersonId); err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred test person's books read in preparation for the test")
	}

	return networkCtx, nil
}


func (test NetworkPartitionTest) Run(network networks.Network, testCtx testsuite.TestContext) {
	// Go doesn't have generics so we have to do this cast first
	castedNetwork := network.(*networks.NetworkContext)

	blockedConnRepartitioner, err := getTwoPartitionsRepartitioner(castedNetwork, true)
	if err != nil {
		testCtx.Fatal(stacktrace.Propagate(err, "An error occurred getting the 2-partition repartitioner with blocked connection"))
	}

	if err != nil {
		testCtx.Fatal(stacktrace.Propagate(err, "An error occurred building the repartitioner block access between API <-> datastore"))
	}

	logrus.Info("Partitioning API and datastore services off from each other...")
	if err := castedNetwork.RepartitionNetwork(blockedConnRepartitioner); err != nil {
		testCtx.Fatal(stacktrace.Propagate(err, "An error occurred repartitioning the network to block access between API <-> datastore"))
	}
	logrus.Info("Repartition complete")

	uncastedApiService, err := castedNetwork.GetService(apiServiceId)
	if err != nil {
		testCtx.Fatal(stacktrace.Propagate(err, "An error occurred getting the API service interface"))
	}
	apiService := uncastedApiService.(*api.ApiService)	// Necessary because Go doesn't have generics

	// Use a short timeout because we expect a partition
	logrus.Info("Incrementing books read while partition is in place, to verify no comms are possible...")
	if err := apiService.IncrementBooksRead(testPersonId); err == nil {
		testCtx.Fatal(stacktrace.NewError("Expected the book increment call to fail due to the network " +
			"partition between API and datastore services, but no error was thrown"))
	} else {
		logrus.Infof("Incrementing books read threw the following error as expected due to network partition: %v", err)
	}

	openConnRepartitioner, err := getTwoPartitionsRepartitioner(castedNetwork, false)
	if err != nil {
		testCtx.Fatal(stacktrace.Propagate(err, "An error occurred getting the 2-partition repartitioner with open connection"))
	}

	logrus.Info("Repartitioning to heal partition between API and datastore...")
	if err := castedNetwork.RepartitionNetwork(openConnRepartitioner); err != nil {
		testCtx.Fatal(stacktrace.Propagate(err, "An error occurred healing the partition"))
	}
	logrus.Info("Partition healed successfully")

	logrus.Info("Making another call to increment books read, where the partition will heal in the middle of the call...")
	// Use infinite timeout because we expect the partition healing to fix the issue
	if err := apiService.IncrementBooksRead(testPersonId); err != nil {
		testCtx.Fatal(stacktrace.Propagate(
			err,
			"An error occurred incrementing the number of books read, even though the partition should have been " +
				"healed by the goroutine",
		))
	}
	logrus.Info("Successfully incremented books read, indicating that the partition has healed successfully!")
}


func (test *NetworkPartitionTest) GetTestConfiguration() testsuite.TestConfiguration {
	return testsuite.TestConfiguration{
		IsPartitioningEnabled: true,
	}
}

func (test NetworkPartitionTest) GetExecutionTimeout() time.Duration {
	return 60 * time.Second
}

func (test NetworkPartitionTest) GetSetupTeardownBuffer() time.Duration {
	return 60 * time.Second
}

/*
Creates a repartitioner that will partition the network between the API & datastore services, with the connection between them configurable
 */
func getTwoPartitionsRepartitioner(
		networkCtx *networks.NetworkContext,
		isConnectionBlocked bool) (*networks.Repartitioner, error){
	repartitioner, err := networkCtx.GetRepartitionerBuilder(
			false,
		).WithPartition(
			apiPartitionId,
			apiServiceId,
		).WithPartition(
			datastorePartitionId,
			datastoreServiceId,
		).WithPartitionConnection(
			apiPartitionId,
			datastorePartitionId,
			isConnectionBlocked,
		).Build()
	if err != nil {
		return nil, stacktrace.Propagate(
			err,
			"An error occurred creating a two-partition repartitioner with isConnectionBlocked = %v",
			isConnectionBlocked)
	}
	return repartitioner, nil
}


