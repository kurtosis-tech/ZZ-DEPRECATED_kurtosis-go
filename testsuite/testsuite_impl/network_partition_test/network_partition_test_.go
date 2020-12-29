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

	gaiaPartitionId = "we-are-all-one"

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


	repartitioner, err := castedNetwork.GetRepartitionerBuilder(
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
			true,
		).Build()
	if err != nil {
		testCtx.Fatal(stacktrace.Propagate(err, "An error occurred building the repartitioner block access between API <-> datastore"))
	}

	if err := castedNetwork.RepartitionNetwork(repartitioner); err != nil {
		testCtx.Fatal(stacktrace.Propagate(err, "An error occurred repartitioning the network to block access between API <-> datastore"))
	}

	uncastedApiService, err := castedNetwork.GetService(apiServiceId)
	if err != nil {
		testCtx.Fatal(stacktrace.Propagate(err, "An error occurred getting the API service interface"))
	}
	apiService := uncastedApiService.(*api.ApiService)	// Necessary because Go doesn't have generics

	if err := apiService.IncrementBooksRead(testPersonId); err == nil {
		testCtx.Fatal(stacktrace.NewError("Expected the book increment call to fail due to the network " +
			"partition between API and datastore services, but no error was thrown"))
	} else {
		logrus.Infof("Incrementing books read threw the following error as expected due to network partition: %v", err)
	}
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


