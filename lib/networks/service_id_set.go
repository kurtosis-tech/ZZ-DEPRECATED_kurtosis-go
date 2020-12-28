/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package networks

import "github.com/kurtosis-tech/kurtosis-go/lib/services"

type serviceIdSet struct {
	elems map[services.ServiceID]bool
}

func newServiceIdSet() *serviceIdSet {
	return &serviceIdSet{
		elems: map[services.ServiceID]bool{},
	}
}

func (set *serviceIdSet) add(id services.ServiceID) {
	set.elems[id] = true
}

func (set *serviceIdSet) contains(id services.ServiceID) bool {
	_, found := set.elems[id]
	return found
}

func (set *serviceIdSet) getElems() []services.ServiceID {
	result := []services.ServiceID{}
	for id, _ := range set.elems {
		result = append(result, id)
	}
	return result
}