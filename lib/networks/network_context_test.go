/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package networks

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/test_mocks/mock_docker_container_initializer"
	"github.com/kurtosis-tech/kurtosis-go/lib/test_mocks/mock_kurtosis_service"
	"testing"
)

func TestDisallowingSameIds(t *testing.T) {
	var duplicatedId ServiceID = "the-id"
	kurtosisService := mock_kurtosis_service.NewMockKurtosisService()
	networkCtx := NewNetworkContext(kurtosisService, "/services")
	_, _, err := networkCtx.AddService(duplicatedId, mock_docker_container_initializer.NewMockDockerContainerInitializer())
	if err != nil {
		t.Fatalf("Expected first service to get added successfully but an error occurred: %v", err)
	}

	_, _, err = networkCtx.AddService(duplicatedId, mock_docker_container_initializer.NewMockDockerContainerInitializer())
	if err == nil {
		t.Fatalf("Expected adding the second service to throw an error due to the duplicate service ID but no error occurred")
	}
}
