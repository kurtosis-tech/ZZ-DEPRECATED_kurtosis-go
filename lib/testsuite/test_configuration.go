/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package testsuite

/*
Holds configuration values that, if set, give the test the ability to do special things
 */
type TestConfiguration struct {
	// If true, enables the test to set up network partitions between services
	// This should NOT be done thoughtlessly, however - when partitioning is enabled,
	//  adding services will be slower because all the other nodes in the network will
	//  need to update their iptables for the new node. The slowdown will scale with the
	//  number of services in your network.
	IsPartitioningEnabled bool
}
