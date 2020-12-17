/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package advanced_network_test

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/networks"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
	"github.com/kurtosis-tech/kurtosis-go/testsuite/networks_impl"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	numApiServices = 2

	testPersonId = 46
)

type AdvancedNetworkTest struct {
	datastoreServiceImage string
	apiServiceImage string
	apiServiceIds []services.ServiceID
}

func NewAdvancedNetworkTest(datastoreServiceImage string, apiServiceImage string) *AdvancedNetworkTest {
	return &AdvancedNetworkTest{datastoreServiceImage: datastoreServiceImage, apiServiceImage: apiServiceImage}
}

func (test *AdvancedNetworkTest) Setup(networkCtx *networks.NetworkContext) (networks.Network, error) {
	network := networks_impl.NewTestNetwork(networkCtx, test.datastoreServiceImage, test.apiServiceImage)

	if err := network.AddDatastore(); err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred adding the datastore")
	}

	for i := 0; i < numApiServices; i++ {
		apiService, err := network.AddApiService()
		if err != nil {
			return nil, stacktrace.Propagate(err, "An error occurred adding API service %v", i)
		}
		test.apiServiceIds = append(test.apiServiceIds, apiService)
	}

	return network, nil
}

func (test *AdvancedNetworkTest) Run(network networks.Network, testCtx testsuite.TestContext) {
	castedNetwork := network.(*networks_impl.TestNetwork)
	firstApiService, err := castedNetwork.GetApiService(test.apiServiceIds[0])
	if err != nil {
		testCtx.Fatal(stacktrace.Propagate(err, "An error occurred getting the first API service"))
	}
	secondApiService, err := castedNetwork.GetApiService(test.apiServiceIds[1])
	if err != nil {
		testCtx.Fatal(stacktrace.Propagate(err, "An error occurred getting the second API service"))
	}

	logrus.Infof("Adding test person via first API service...")
	if err := firstApiService.AddPerson(testPersonId); err != nil {
		testCtx.Fatal(stacktrace.Propagate(err, "An error occurred adding test person through first API service"))
	}
	logrus.Info("Test person added")

	logrus.Infof("Incrementing test person's number of books read through first API service ...")
	if err := firstApiService.IncrementBooksRead(testPersonId); err != nil {
		testCtx.Fatal(stacktrace.Propagate(err, "An error occurred incrementing the number of books read through the first API service"))
	}
	logrus.Info("Incremented number of books read")

	logrus.Info("Retrieving test person to verify number of books read through second API service...")
	person, err := secondApiService.GetPerson(testPersonId)
	if err != nil {
		testCtx.Fatal(stacktrace.Propagate(err, "An error occurred getting the test person through second API service"))
	}
	logrus.Info("Retrieved test person through second API service")

	testCtx.AssertTrue(
		person.BooksRead == 1,
		stacktrace.NewError(
			"Expected number of books read to be incremented, but was '%v'",
			person.BooksRead,
		),
	)
}

func (test AdvancedNetworkTest) GetExecutionTimeout() time.Duration {
	return 60 * time.Second
}

func (test AdvancedNetworkTest) GetSetupTeardownBuffer() time.Duration {
	return 60 * time.Second
}




