package ssz

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type sszSizeHint struct {
	size    uint64
	dynamic bool
	specval bool
}

func (d *DynSsz) getSszSizeTag(field *reflect.StructField) ([]sszSizeHint, error) {
	sszSizes := []sszSizeHint{}

	if fieldSszSizeStr, fieldHasSszSize := field.Tag.Lookup("ssz-size"); fieldHasSszSize {
		for _, sszSizeStr := range strings.Split(fieldSszSizeStr, ",") {
			sszSize := sszSizeHint{}

			if sszSizeStr == "?" {
				sszSize.dynamic = true
			} else {
				sszSizeInt, err := strconv.ParseUint(sszSizeStr, 10, 32)
				if err != nil {
					return sszSizes, fmt.Errorf("error parsing ssz-size tag for '%v' field: %v", field.Name, err)
				}
				sszSize.size = sszSizeInt
			}

			sszSizes = append(sszSizes, sszSize)
		}
	}

	fieldDynSszSizeStr, fieldHasDynSszSize := field.Tag.Lookup("dynssz-size")
	if fieldHasDynSszSize {
		for i, sszSizeStr := range strings.Split(fieldDynSszSizeStr, ",") {
			sszSize := sszSizeHint{}

			if sszSizeStr == "?" {
				sszSize.dynamic = true
			} else if sszSizeInt, err := strconv.ParseUint(sszSizeStr, 10, 32); err == nil {
				sszSize.size = sszSizeInt
			} else {
				ok, specVal, err := d.getSpecValue(sszSizeStr)
				if err != nil {
					return sszSizes, fmt.Errorf("error parsing dynssz-size tag for '%v' field (%v): %v", field.Name, sszSizeStr, err)
				}
				if ok {
					// dynamic value from spec
					sszSize.size = specVal
					sszSize.specval = true
				} else {
					// unknown spec value? fallback to fastssz
					break
				}
			}

			if sszSizes[i].size != sszSize.size {
				sszSizes[i] = sszSize
			}
		}
	}

	return sszSizes, nil
}

func (d *DynSsz) getSszSize(targetType reflect.Type, sizeHints []sszSizeHint) (int, bool, error) {
	staticSize := 0
	hasSpecValue := false
	isDynamicSize := false

	childSizeHints := []sszSizeHint{}
	if len(sizeHints) > 1 {
		childSizeHints = sizeHints[1:]
	}

	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}

	switch targetType.Kind() {
	case reflect.Struct:
		for i := 0; i < targetType.NumField(); i++ {
			field := targetType.Field(i)
			size, hasSpecVal, _, err := d.getSszFieldSize(&field)
			if err != nil {
				return 0, false, err
			}
			if size < 0 {
				isDynamicSize = true
			}
			if hasSpecVal {
				hasSpecValue = true
			}
			staticSize += size
		}
	case reflect.Array:
		arrLen := targetType.Len()
		fieldType := targetType.Elem()
		size, hasSpecVal, err := d.getSszSize(fieldType, childSizeHints)
		if err != nil {
			return 0, false, err
		}
		if size < 0 {
			isDynamicSize = true
		}
		if hasSpecVal {
			hasSpecValue = true
		}
		staticSize += size * arrLen
	case reflect.Slice:
		fieldType := targetType.Elem()
		size, hasSpecVal, err := d.getSszSize(fieldType, childSizeHints)
		if err != nil {
			return 0, false, err
		}
		if size < 0 {
			isDynamicSize = true
		}
		if hasSpecVal || (len(sizeHints) > 0 && sizeHints[0].specval) {
			hasSpecValue = true
		}

		if len(sizeHints) > 0 && sizeHints[0].size > 0 {
			staticSize += size * int(sizeHints[0].size)
		} else {
			isDynamicSize = true
		}
	case reflect.Bool:
		staticSize = 1
	case reflect.Uint8:
		staticSize = 1
	case reflect.Uint16:
		staticSize = 2
	case reflect.Uint32:
		staticSize = 4
	case reflect.Uint64:
		staticSize = 8
	default:
		return 0, false, fmt.Errorf("unhandled reflection kind in size check: %v", targetType.Kind())
	}

	if isDynamicSize {
		staticSize = -1
	}

	return staticSize, hasSpecValue, nil
}

func (d *DynSsz) getSszFieldSize(targetField *reflect.StructField) (int, bool, []sszSizeHint, error) {
	sszSizes, err := d.getSszSizeTag(targetField)
	if err != nil {
		return 0, false, nil, err
	}

	size, hasSpecVal, err := d.getSszSize(targetField.Type, sszSizes)
	return size, hasSpecVal, sszSizes, err
}

func (d *DynSsz) getSszValueSize(targetType reflect.Type, targetValue reflect.Value) (int, error) {
	staticSize := 0

	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
		targetValue = targetValue.Elem()
	}

	switch targetType.Kind() {
	case reflect.Struct:
		for i := 0; i < targetType.NumField(); i++ {
			field := targetType.Field(i)
			fieldValue := targetValue.Field(i)

			fieldTypeSize, _, _, err := d.getSszFieldSize(&field)
			if err != nil {
				return 0, err
			}

			if fieldTypeSize < 0 {
				// dynamic field, add 4 bytes for offset
				staticSize += 4
			}

			size, err := d.getSszValueSize(field.Type, fieldValue)
			if err != nil {
				return 0, err
			}

			staticSize += size
		}
	case reflect.Array:
		arrLen := targetType.Len()
		if arrLen > 0 {
			fieldType := targetType.Elem()
			size, err := d.getSszValueSize(fieldType, targetValue.Index(0))
			if err != nil {
				return 0, err
			}
			staticSize = size * arrLen
		}
	case reflect.Slice:
		fieldType := targetType.Elem()
		sliceLen := targetValue.Len()

		if sliceLen > 0 {
			size, err := d.getSszValueSize(fieldType, targetValue.Index(0))
			if err != nil {
				return 0, err
			}
			staticSize = size * sliceLen
		}
	case reflect.Bool:
		staticSize = 1
	case reflect.Uint8:
		staticSize = 1
	case reflect.Uint16:
		staticSize = 2
	case reflect.Uint32:
		staticSize = 4
	case reflect.Uint64:
		staticSize = 8
	default:
		return 0, fmt.Errorf("unhandled reflection kind in size check: %v", targetType.Kind())
	}

	return staticSize, nil
}
