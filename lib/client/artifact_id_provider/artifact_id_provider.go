/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package artifact_id_provider

import (
	"crypto"
	"encoding/hex"
	"github.com/palantir/stacktrace"

	// This is a special type of import that includes the correct hashing algorithm that we use
	// If we don't have the "_" in front, Goland will complain it's unused
	_ "golang.org/x/crypto/sha3"
)

const (
	defaultHashFunction = crypto.SHA3_256
)

type ArtifactIdProvider interface {
	GetArtifactId(artifactUrl string) (id string, resultErr error)
}

type DefaultArtifactIdProvider struct {
	hashFunction crypto.Hash
}

func NewDefaultArtifactIdProvider() *DefaultArtifactIdProvider {
	return &DefaultArtifactIdProvider{hashFunction: defaultHashFunction}
}

// Gets a unique ID for an artifact as identified by its URL
func (defaultProvider DefaultArtifactIdProvider) GetArtifactId(artifactUrl string) (id string, resultErr error) {
	hasher := defaultHashFunction.New()
	artifactUrlBytes := []byte(artifactUrl)
	if _, err := hasher.Write(artifactUrlBytes); err != nil {
		return "", stacktrace.Propagate(err, "An error occurred writing the artifact URL to the hash function")
	}
	hexEncodedHash := hex.EncodeToString(hasher.Sum(nil))
	return hexEncodedHash, nil
}
