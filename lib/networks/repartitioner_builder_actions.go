/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package networks

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/palantir/stacktrace"
)

type repartitionerMutator interface {
	mutate(repartitioner *Repartitioner) error
}

// ======================================================================================================
//                                         Add Partition
// ======================================================================================================
type addPartitionAction struct {
	partition PartitionID
	services  []services.ServiceID
}

func (a addPartitionAction) mutate(repartitioner *Repartitioner) error {
	newPartition := a.partition
	_, found := repartitioner.partitionServices[newPartition]
	if found {
		return stacktrace.NewError("Partition '%v' already declared",newPartition)
	}
	newPartitionServices := newServiceIdSet()
	for _, id := range a.services {
		newPartitionServices.add(id)
	}
	repartitioner.partitionServices[newPartition] = newPartitionServices
	return nil
}

// ======================================================================================================
//                                         Add Partition
// ======================================================================================================
type addPartitionConnectionAction struct {
	partitionA PartitionID
	partitionB PartitionID
	connection PartitionConnection
}

func (a addPartitionConnectionAction) mutate(repartitioner *Repartitioner) error {
	partitionA := a.partitionA
	partitionB := a.partitionB
	connectionInfo := a.connection

	// The API service will already check the forward & reverse, so we can limit our error-checking here to ensuring the user
	//  doesn't overwrite something they've already defined
	partitionAConns, foundA := repartitioner.partitionConnections[partitionA]
	if foundA {
		if _, foundB := partitionAConns[partitionB]; foundB {
			return stacktrace.NewError(
				"Partition connection '%v' <-> '%v' is already defined",
				partitionA,
				partitionB)
		}
	}

	partitionAConns, found := repartitioner.partitionConnections[partitionA]
	if !found {
		partitionAConns = map[PartitionID]PartitionConnection{}
	}
	partitionAConns[partitionB] = connectionInfo
	repartitioner.partitionConnections[partitionA] = partitionAConns
	return nil
}
