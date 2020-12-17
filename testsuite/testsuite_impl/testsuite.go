/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package testsuite_impl

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/testsuite"
	"github.com/kurtosis-tech/kurtosis-go/testsuite/testsuite_impl/advanced_network_test"
	"github.com/kurtosis-tech/kurtosis-go/testsuite/testsuite_impl/basic_datastore_and_api_test"
	"github.com/kurtosis-tech/kurtosis-go/testsuite/testsuite_impl/basic_datastore_test"
)

type Testsuite struct {
	apiServiceImage string
	datastoreServiceImage string
}

func NewTestsuite(apiServiceImage string, datastoreServiceImage string) *Testsuite {
	return &Testsuite{apiServiceImage: apiServiceImage, datastoreServiceImage: datastoreServiceImage}
}

func (suite Testsuite) GetTests() map[string]testsuite.Test {
	return map[string]testsuite.Test{
		"basicDatastoreTest": basic_datastore_test.NewBasicDatastoreTest(suite.datastoreServiceImage),
		"basicDatastoreAndApiTest": basic_datastore_and_api_test.NewBasicDatastoreAndApiTest(
			suite.datastoreServiceImage,
			suite.apiServiceImage,
		),
		"advancedNetworkTest": advanced_network_test.NewAdvancedNetworkTest(
			suite.datastoreServiceImage,
			suite.apiServiceImage,
		),
	}
}

func (suite Testsuite) GetNetworkWidthBits() uint32 {
	return 8
}


