package ssz

import "reflect"

type DynSsz struct {
	typesWithSpecVals map[reflect.Type]uint8
	SpecValues        map[string]any
}

func NewDynSsz() *DynSsz {
	return &DynSsz{
		typesWithSpecVals: map[reflect.Type]uint8{},
		SpecValues: map[string]any{
			"SLOTS_PER_HISTORICAL_ROOT": uint64(8192),
		},
	}
}
