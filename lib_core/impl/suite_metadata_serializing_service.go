/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package impl

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-go/lib_core/api/generated"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SuiteMetadataSerializingService struct {}

func (s SuiteMetadataSerializingService) SerializeSuiteMetadata(
		ctx context.Context,
		metadata *generated.TestSuiteMetadata) (*emptypb.Empty, error) {

	return nil, nil
}

