/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package networks

type PartitionID string

type PartitionConnection struct {
	IsBlocked bool
}

type Repartitioner struct {
	partitionServices map[PartitionID]*serviceIdSet
	partitionConnections map[PartitionID]map[PartitionID]PartitionConnection
	defaultConnection PartitionConnection
}
