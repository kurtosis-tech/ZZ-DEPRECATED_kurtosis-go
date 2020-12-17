/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package networks

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/kurtosis-tech/kurtosis-go/lib/test_mocks/mock_docker_container_initializer"
	"github.com/kurtosis-tech/kurtosis-go/lib/test_mocks/mock_kurtosis_service"
	"io/ioutil"
	"testing"
)

func TestDisallowingSameIds(t *testing.T) {
	var duplicatedId services.ServiceID = "the-id"
	kurtosisService := mock_kurtosis_service.NewMockKurtosisService()

	tempDirpath, err := ioutil.TempDir("", "suite-volume")
	if err != nil {
		t.Fatalf("An error occurred creating the temporary directory to represent the suite execution volume: %v", err)
	}

	networkCtx := NewNetworkContext(kurtosisService, tempDirpath, "/")
	_, _, err = networkCtx.AddService(duplicatedId, mock_docker_container_initializer.NewMockDockerContainerInitializer())
	if err != nil {
		t.Fatalf("Expected first service to get added successfully but an error occurred: %v", err)
	}

	_, _, err = networkCtx.AddService(duplicatedId, mock_docker_container_initializer.NewMockDockerContainerInitializer())
	if err == nil {
		t.Fatalf("Expected adding the second service to throw an error due to the duplicate service ID but no error occurred")
	}
}
