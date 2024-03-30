package ssz

import (
	"reflect"

	fastssz "github.com/ferranbt/fastssz"
)

var byteType = reflect.TypeOf(byte(0))
var sszMarshallerType = reflect.TypeOf((*fastssz.Marshaler)(nil)).Elem()
var sszUnmarshallerType = reflect.TypeOf((*fastssz.Unmarshaler)(nil)).Elem()
