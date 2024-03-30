package ssz

import (
	"fmt"
	"reflect"

	fastssz "github.com/ferranbt/fastssz"
)

func UnmarshalSSZ(target any, ssz []byte) error {
	d := NewDynSsz()

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

func (d *DynSsz) unmarshalType(targetType reflect.Type, targetValue reflect.Value, ssz []byte, sizeHints []sszSizeHint, idt int) (int, error) {
	consumedBytes := 0

	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
		if targetValue.IsNil() {
			//fmt.Printf("new type %v\n", targetType.Name())
			newValue := reflect.New(targetType)
			targetValue.Set(newValue)
		}
		targetValue = targetValue.Elem()
	}

	//fmt.Printf("%stype: %s\t kind: %v\n", indent(idt), targetType.Name(), targetType.Kind())

	switch targetType.Kind() {
	case reflect.Struct:
		usedFastSsz := false

		hasSpecVals := d.typesWithSpecVals[targetType]
		if hasSpecVals == unknownSpecValued {
			hasSpecVals = noSpecValues
			if targetValue.Addr().Type().Implements(sszUnmarshallerType) {
				_, hasSpecVals2, err := d.getSszSize(targetType, sizeHints)
				if err != nil {
					return 0, err
				}

				if hasSpecVals2 {
					hasSpecVals = hasSpecValues
				}
			}

			//fmt.Printf("%s fastssz for type %s: %v\n", indent(idt), targetType.Name(), hasSpecVals)
			d.typesWithSpecVals[targetType] = hasSpecVals
		}
		if hasSpecVals == noSpecValues {
			unmarshaller, ok := targetValue.Addr().Interface().(fastssz.Unmarshaler)
			if ok {
				err := unmarshaller.UnmarshalSSZ(ssz)
				if err != nil {
					return 0, err
				}
				consumedBytes = len(ssz)
				usedFastSsz = true
			}
		}

		if !usedFastSsz {
			consumed, err := d.unmarshalStruct(targetType, targetValue, ssz, idt)
			if err != nil {
				return 0, err
			}
			consumedBytes = consumed
		}
	case reflect.Array:
		consumed, err := d.unmarshalArray(targetType, targetValue, ssz, sizeHints, idt)
		if err != nil {
			return 0, err
		}
		consumedBytes = consumed
	case reflect.Slice:
		consumed, err := d.unmarshalSlice(targetType, targetValue, ssz, sizeHints, idt)
		if err != nil {
			return 0, err
		}
		consumedBytes = consumed
	case reflect.Bool:
		targetValue.SetBool(fastssz.UnmarshalBool(ssz))
		consumedBytes = 1
	case reflect.Uint8:
		targetValue.SetUint(uint64(fastssz.UnmarshallUint8(ssz)))
		consumedBytes = 1
	case reflect.Uint16:
		targetValue.SetUint(uint64(fastssz.UnmarshallUint16(ssz)))
		consumedBytes = 2
	case reflect.Uint32:
		targetValue.SetUint(uint64(fastssz.UnmarshallUint32(ssz)))
		consumedBytes = 4
	case reflect.Uint64:
		targetValue.SetUint(uint64(fastssz.UnmarshallUint64(ssz)))
		consumedBytes = 8
	default:
		return 0, fmt.Errorf("unknown type: %v", targetType)
	}

	return consumedBytes, nil
}

func (d *DynSsz) unmarshalStruct(targetType reflect.Type, targetValue reflect.Value, ssz []byte, idt int) (int, error) {
	offset := 0
	dynamicFields := []*reflect.StructField{}
	dynamicOffsets := []int{}
	dynamicSizeHints := [][]sszSizeHint{}

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)

		fieldSize, _, sizeHints, err := d.getSszFieldSize(&field)
		if err != nil {
			return 0, err
		}

		if fieldSize > 0 {
			//fmt.Printf("%sfield %d:\t static [%v:%v] %v\t %v\n", indent(idt+1), i, offset, offset+fieldSize, fieldSize, field.Name)

			fieldSsz := ssz[offset : offset+fieldSize]
			fieldValue := targetValue.Field(i)
			consumedBytes, err := d.unmarshalType(field.Type, fieldValue, fieldSsz, sizeHints, idt+2)
			if err != nil {
				return 0, fmt.Errorf("failed decoding field %v: %v", field.Name, err)
			}
			if consumedBytes != fieldSize {
				return 0, fmt.Errorf("struct field did not consume expected ssz range (consumed: %v, expected: %v)", consumedBytes, fieldSize)
			}

		} else {
			fieldSize = 4
			fieldOffset := fastssz.ReadOffset(ssz[offset : offset+fieldSize])
			//fmt.Printf("%sfield %d:\t offset [%v:%v] %v\t %v \t %v\n", indent(idt+1), i, offset, offset+fieldSize, fieldSize, field.Name, fieldOffset)

			dynamicFields = append(dynamicFields, &field)
			dynamicOffsets = append(dynamicOffsets, int(fieldOffset))
			dynamicSizeHints = append(dynamicSizeHints, sizeHints)
		}
		offset += fieldSize
	}
	dynamicFieldCount := len(dynamicFields)
	for i, field := range dynamicFields {
		var endOffset int
		startOffset := dynamicOffsets[i]
		if i < dynamicFieldCount-1 {
			endOffset = dynamicOffsets[i+1]
		} else {
			endOffset = len(ssz)
		}

		//fmt.Printf("%sfield %d:\t dynamic [%v:%v]\t %v\n", indent(idt+1), field.Index[0], startOffset, endOffset, field.Name)

		var fieldSsz []byte
		if endOffset > startOffset {
			fieldSsz = ssz[startOffset:endOffset]
		} else {
			fieldSsz = []byte{}
		}

		fieldValue := targetValue.Field(field.Index[0])
		consumedBytes, err := d.unmarshalType(field.Type, fieldValue, fieldSsz, dynamicSizeHints[i], idt+2)
		if err != nil {
			return 0, fmt.Errorf("failed decoding field %v: %v", field.Name, err)
		}
		if consumedBytes != endOffset-startOffset {
			return 0, fmt.Errorf("struct field did not consume expected ssz range (consumed: %v, expected: %v)", consumedBytes, endOffset-startOffset)
		}

		offset += consumedBytes
	}

	return offset, nil
}

func (d *DynSsz) unmarshalArray(targetType reflect.Type, targetValue reflect.Value, ssz []byte, sizeHints []sszSizeHint, idt int) (int, error) {
	var consumedBytes int

	childSizeHints := []sszSizeHint{}
	if len(sizeHints) > 1 {
		childSizeHints = sizeHints[1:]
	}

	fieldType := targetType.Elem()
	fieldIsPtr := fieldType.Kind() == reflect.Ptr
	if fieldIsPtr {
		fieldType = fieldType.Elem()
	}

	arrLen := targetType.Len()
	if fieldType == byteType {
		// shortcut for performance: use copy on []byte arrays
		reflect.Copy(targetValue, reflect.ValueOf(ssz[0:arrLen]))
		consumedBytes = arrLen
	} else {
		offset := 0
		itemSize := len(ssz) / arrLen
		for i := 0; i < arrLen; i++ {
			var itemVal reflect.Value
			if fieldIsPtr {
				//fmt.Printf("new array item %v\n", fieldType.Name())
				itemVal = reflect.New(fieldType).Elem()
				targetValue.Index(i).Set(itemVal.Addr())
			} else {
				itemVal = targetValue.Index(i)
			}

			itemSsz := ssz[offset : offset+itemSize]

			consumed, err := d.unmarshalType(fieldType, itemVal, itemSsz, childSizeHints, idt+2)
			if err != nil {
				return 0, err
			}
			if consumed != itemSize {
				return 0, fmt.Errorf("unmarshalling array item did not consume expected ssz range (consumed: %v, expected: %v)", consumed, itemSize)
			}

			offset += itemSize
		}

		consumedBytes = offset
	}

	return consumedBytes, nil
}

func (d *DynSsz) unmarshalSlice(targetType reflect.Type, targetValue reflect.Value, ssz []byte, sizeHints []sszSizeHint, idt int) (int, error) {
	var consumedBytes int

	childSizeHints := []sszSizeHint{}
	if len(sizeHints) > 1 {
		childSizeHints = sizeHints[1:]
	}

	fieldType := targetType.Elem()
	fieldIsPtr := fieldType.Kind() == reflect.Ptr
	if fieldIsPtr {
		fieldType = fieldType.Elem()
	}

	sliceLen := 0
	sszLen := len(ssz)

	if len(sizeHints) > 0 && sizeHints[0].size > 0 {
		sliceLen = int(sizeHints[0].size)
	} else if len(sizeHints) > 1 && sizeHints[1].size > 0 {
		ok := false
		sliceLen, ok = fastssz.DivideInt(sszLen, int(sizeHints[1].size))
		if !ok {
			return 0, fmt.Errorf("invalid slice length, expected multiple of %v, got %v", sizeHints[1], sszLen)
		}
	} else {
		size, _, err := d.getSszSize(fieldType, childSizeHints)
		if err != nil {
			return 0, err
		}

		if size > 0 {
			ok := false
			sliceLen, ok = fastssz.DivideInt(sszLen, size)
			if !ok {
				return 0, fmt.Errorf("invalid slice length, expected multiple of %v, got %v", size, sszLen)
			}
		}
	}

	if sliceLen == 0 && len(ssz) > 0 {
		return 0, fmt.Errorf("cannot deteriminate length of dynamic slice")
	}

	//fmt.Printf("new slice %v  %v\n", fieldType.Name(), sliceLen)
	newValue := reflect.MakeSlice(targetType, sliceLen, sliceLen)
	targetValue.Set(newValue)

	if fieldType == byteType {
		// shortcut for performance: use copy on []byte arrays
		reflect.Copy(newValue, reflect.ValueOf(ssz[0:sliceLen]))
		consumedBytes = sliceLen
	} else {
		offset := 0
		if sliceLen > 0 {
			itemSize := sszLen / sliceLen

			for i := 0; i < sliceLen; i++ {
				var itemVal reflect.Value
				if fieldIsPtr {
					//fmt.Printf("new slice item %v\n", fieldType.Name())
					itemVal = reflect.New(fieldType).Elem()
					newValue.Index(i).Set(itemVal.Addr())
				} else {
					itemVal = newValue.Index(i)
				}

				itemSsz := ssz[offset : offset+itemSize]

				consumed, err := d.unmarshalType(fieldType, itemVal, itemSsz, childSizeHints, idt+2)
				if err != nil {
					return 0, err
				}
				if consumed != itemSize {
					return 0, fmt.Errorf("slice item did not consume expected ssz range (consumed: %v, expected: %v)", consumed, itemSize)
				}

				offset += itemSize
			}
		}

		consumedBytes = offset
	}

	return consumedBytes, nil
}
