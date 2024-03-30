package ssz

import (
	"fmt"
	"reflect"
)

type DynSsz struct {
	typesWithSpecVals map[reflect.Type]uint8
	SpecValues        map[string]any
	NoFastSsz         bool
}

const (
	unknownSpecValued uint8 = iota
	noSpecValues
	hasSpecValues
)

func NewDynSsz(specs map[string]any) *DynSsz {
	if specs == nil {
		specs = map[string]any{}
	}
	return &DynSsz{
		typesWithSpecVals: map[reflect.Type]uint8{},
		SpecValues:        specs,
	}
}

func (d *DynSsz) MarshalSSZ(source any) ([]byte, error) {
	sourceType := reflect.TypeOf(source)
	sourceValue := reflect.ValueOf(source)

	size, err := d.getSszValueSize(sourceType, sourceValue)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, size)
	newBuf, err := d.marshalType(sourceType, sourceValue, buf[:0], []sszSizeHint{}, 0)
	if err != nil {
		return nil, err
	}

	if len(newBuf) != size {
		return nil, fmt.Errorf("ssz length does not match expected length (expected: %v, got: %v)", size, len(newBuf))
	}

	return newBuf, nil
}

func (d *DynSsz) MarshalSSZTo(source any, buf []byte) ([]byte, error) {
	sourceType := reflect.TypeOf(source)
	sourceValue := reflect.ValueOf(source)

	newBuf, err := d.marshalType(sourceType, sourceValue, buf, []sszSizeHint{}, 0)
	if err != nil {
		return nil, err
	}

	return newBuf, nil
}

func (d *DynSsz) SizeSSZ(source any) (int, error) {
	sourceType := reflect.TypeOf(source)
	sourceValue := reflect.ValueOf(source)

	return d.getSszValueSize(sourceType, sourceValue)
}

func (d *DynSsz) UnmarshalSSZ(target any, ssz []byte) error {
	targetType := reflect.TypeOf(target)
	targetValue := reflect.ValueOf(target)

	consumedBytes, err := d.unmarshalType(targetType, targetValue, ssz, []sszSizeHint{}, 0)
	if err != nil {
		return err
	}

	if consumedBytes != len(ssz) {
		return fmt.Errorf("did not consume full ssz range (consumed: %v, ssz size: %v)", consumedBytes, len(ssz))
	}

	return nil
}
