/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package method_types

type SerializablePartitionConnection struct {
	IsBlocked bool		`json:"isBlocked"`
}

type RepartitionArgs struct {
	// Mapping of partition ID -> "set" of service IDs
	PartitionServices map[string]map[string]bool	`json:"partitionServices"`

	// Mapping of partitionA -> partitionB -> partitionConnection details
	// We use this format because JSON doesn't allow object keys
	PartitionConnections map[string]map[string]SerializablePartitionConnection `json:"partitionConnections"`

	DefaultConnection SerializablePartitionConnection `json:"defaultConnection"`
}
