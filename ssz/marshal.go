package ssz

import (
	"encoding/binary"
	"fmt"
	"reflect"

	fastssz "github.com/ferranbt/fastssz"
)

func MarshalSSZ(source any, buf []byte) ([]byte, error) {
	d := NewDynSsz()
	d.NoFastSsz = true

	sourceType := reflect.TypeOf(source)
	sourceValue := reflect.ValueOf(source)

	newBuf, err := d.marshalType(sourceType, sourceValue, buf, []sszSizeHint{}, 0)
	if err != nil {
		return nil, err
	}

	return newBuf, nil
}

func (d *DynSsz) marshalType(sourceType reflect.Type, sourceValue reflect.Value, buf []byte, sizeHints []sszSizeHint, idt int) ([]byte, error) {
	if sourceType.Kind() == reflect.Ptr {
		sourceType = sourceType.Elem()
		sourceValue = sourceValue.Elem()
	}

	//fmt.Printf("%stype: %s\t kind: %v\n", indent(idt), sourceType.Name(), sourceType.Kind())

	switch sourceType.Kind() {
	case reflect.Struct:
		usedFastSsz := false

		hasSpecVals := d.typesWithSpecVals[sourceType]
		if hasSpecVals == unknownSpecValued && !d.NoFastSsz {
			hasSpecVals = noSpecValues
			if sourceValue.Addr().Type().Implements(sszMarshallerType) {
				_, hasSpecVals2, err := d.getSszSize(sourceType, sizeHints)
				if err != nil {
					return nil, err
				}

				if hasSpecVals2 {
					hasSpecVals = hasSpecValues
				}
			}

			fmt.Printf("%s fastssz for type %s: %v\n", indent(idt), sourceType.Name(), hasSpecVals)
			d.typesWithSpecVals[sourceType] = hasSpecVals
		}
		if hasSpecVals == noSpecValues && !d.NoFastSsz {
			marshaller, ok := sourceValue.Addr().Interface().(fastssz.Marshaler)
			if ok {
				newBuf, err := marshaller.MarshalSSZTo(buf)
				if err != nil {
					return nil, err
				}
				buf = newBuf
				usedFastSsz = true
			}
		}

		if !usedFastSsz {
			newBuf, err := d.marshalStruct(sourceType, sourceValue, buf, idt)
			if err != nil {
				return nil, err
			}
			buf = newBuf
		}
	case reflect.Array:
		newBuf, err := d.marshalArray(sourceType, sourceValue, buf, sizeHints, idt)
		if err != nil {
			return nil, err
		}
		buf = newBuf
	case reflect.Slice:
		newBuf, err := d.marshalSlice(sourceType, sourceValue, buf, sizeHints, idt)
		if err != nil {
			return nil, err
		}
		buf = newBuf
	case reflect.Bool:
		buf = fastssz.MarshalBool(buf, sourceValue.Bool())
	case reflect.Uint8:
		buf = fastssz.MarshalUint8(buf, uint8(sourceValue.Uint()))
	case reflect.Uint16:
		buf = fastssz.MarshalUint16(buf, uint16(sourceValue.Uint()))
	case reflect.Uint32:
		buf = fastssz.MarshalUint32(buf, uint32(sourceValue.Uint()))
	case reflect.Uint64:
		buf = fastssz.MarshalUint64(buf, uint64(sourceValue.Uint()))
	default:
		return nil, fmt.Errorf("unknown type: %v", sourceType)
	}

	return buf, nil
}

func (d *DynSsz) marshalStruct(sourceType reflect.Type, sourceValue reflect.Value, buf []byte, idt int) ([]byte, error) {
	offset := 0
	startLen := len(buf)
	dynamicFields := []*reflect.StructField{}
	dynamicOffsets := []int{}
	dynamicSizeHints := [][]sszSizeHint{}

	for i := 0; i < sourceType.NumField(); i++ {
		field := sourceType.Field(i)

		fieldSize, _, sizeHints, err := d.getSszFieldSize(&field)
		if err != nil {
			return nil, err
		}

		if fieldSize > 0 {
			//fmt.Printf("%sfield %d:\t static [%v:%v] %v\t %v\n", indent(idt+1), i, offset, offset+fieldSize, fieldSize, field.Name)

			fieldValue := sourceValue.Field(i)
			newBuf, err := d.marshalType(field.Type, fieldValue, buf, sizeHints, idt+2)
			if err != nil {
				return nil, fmt.Errorf("failed encoding field %v: %v", field.Name, err)
			}
			buf = newBuf

		} else {
			fieldSize = 4
			buf = append(buf, 0, 0, 0, 0)
			//fmt.Printf("%sfield %d:\t offset [%v:%v] %v\t %v\n", indent(idt+1), i, offset, offset+fieldSize, fieldSize, field.Name)

			dynamicFields = append(dynamicFields, &field)
			dynamicOffsets = append(dynamicOffsets, offset)
			dynamicSizeHints = append(dynamicSizeHints, sizeHints)
		}
		offset += fieldSize
	}

	for i, field := range dynamicFields {
		// set field offset
		fieldOffset := dynamicOffsets[i]
		offsetBuf := make([]byte, 4)
		binary.LittleEndian.PutUint32(offsetBuf, uint32(offset))
		copy(buf[fieldOffset+startLen:fieldOffset+startLen+4], offsetBuf)

		//fmt.Printf("%sfield %d:\t dynamic [%v:]\t %v\n", indent(idt+1), field.Index[0], offset, field.Name)

		fieldValue := sourceValue.Field(field.Index[0])
		bufLen := len(buf)
		newBuf, err := d.marshalType(field.Type, fieldValue, buf, dynamicSizeHints[i], idt+2)
		if err != nil {
			return nil, fmt.Errorf("failed decoding field %v: %v", field.Name, err)
		}
		buf = newBuf
		offset += len(buf) - bufLen
	}

	return buf, nil
}

func (d *DynSsz) marshalArray(sourceType reflect.Type, sourceValue reflect.Value, buf []byte, sizeHints []sszSizeHint, idt int) ([]byte, error) {

	childSizeHints := []sszSizeHint{}
	if len(sizeHints) > 1 {
		childSizeHints = sizeHints[1:]
	}

	fieldType := sourceType.Elem()
	fieldIsPtr := fieldType.Kind() == reflect.Ptr
	if fieldIsPtr {
		fieldType = fieldType.Elem()
	}

	arrLen := sourceType.Len()
	if fieldType == byteType {
		// shortcut for performance: use append on []byte arrays
		buf = append(buf, sourceValue.Bytes()...)
	} else {
		for i := 0; i < arrLen; i++ {
			itemVal := sourceValue.Index(i)
			if fieldIsPtr {
				itemVal = itemVal.Elem()
			}

			newBuf, err := d.marshalType(fieldType, itemVal, buf, childSizeHints, idt+2)
			if err != nil {
				return nil, err
			}
			buf = newBuf
		}
	}

	return buf, nil
}

func (d *DynSsz) marshalSlice(sourceType reflect.Type, sourceValue reflect.Value, buf []byte, sizeHints []sszSizeHint, idt int) ([]byte, error) {
	childSizeHints := []sszSizeHint{}
	if len(sizeHints) > 1 {
		childSizeHints = sizeHints[1:]
	}

	fieldType := sourceType.Elem()
	fieldIsPtr := fieldType.Kind() == reflect.Ptr
	if fieldIsPtr {
		fieldType = fieldType.Elem()
	}

	sliceLen := sourceValue.Len()
	if fieldType == byteType {
		// shortcut for performance: use append on []byte arrays
		buf = append(buf, sourceValue.Bytes()...)
	} else {
		for i := 0; i < sliceLen; i++ {
			itemVal := sourceValue.Index(i)
			if fieldIsPtr {
				itemVal = itemVal.Elem()
			}

			newBuf, err := d.marshalType(fieldType, itemVal, buf, childSizeHints, idt+2)
			if err != nil {
				return nil, err
			}
			buf = newBuf
		}
	}

	return buf, nil
}
