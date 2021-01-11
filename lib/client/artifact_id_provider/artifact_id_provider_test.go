/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package artifact_id_provider

import (
	"github.com/palantir/stacktrace"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetArtifactId(t *testing.T) {
	input := "https://www.google.com"
	artifactProvider := NewDefaultArtifactIdProvider()
	hexEncodedHashStr, err := artifactProvider.GetArtifactId(input)
	if err != nil {
		t.Fatal(stacktrace.Propagate(err, "Received an error when hashing artifact URL"))
	}
	expected := "23ac8f7b65bce49bdd0a9a24bebeb4d347a839153315c01cbc8a7bf6f0c8f083"
	assert.Equal(t, expected, string(hexEncodedHashStr))
}
