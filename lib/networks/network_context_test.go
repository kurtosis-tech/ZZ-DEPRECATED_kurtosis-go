/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package networks

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/client/artifact_id_provider"
	"github.com/kurtosis-tech/kurtosis-go/lib/kurtosis_service"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"io/ioutil"
	"testing"
)

func TestDisallowingSameIds(t *testing.T) {
	var duplicatedId services.ServiceID = "the-id"
	kurtosisService := kurtosis_service.NewMockKurtosisService()

	tempDirpath, err := ioutil.TempDir("", "suite-volume")
	if err != nil {
		t.Fatalf("An error occurred creating the temporary directory to represent the suite execution volume: %v", err)
	}

	networkCtx := NewNetworkContext(kurtosisService, tempDirpath, "/", map[services.FilesArtifactID]artifact_id_provider.ArtifactID{})
	_, _, err = networkCtx.AddService(duplicatedId, services.NewMockDockerContainerInitializer())
	if err != nil {
		t.Fatalf("Expected first service to get added successfully but an error occurred: %v", err)
	}

	_, _, err = networkCtx.AddService(duplicatedId, services.NewMockDockerContainerInitializer())
	if err == nil {
		t.Fatalf("Expected adding the second service to throw an error due to the duplicate service ID but no error occurred")
	}
}
