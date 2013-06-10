package mirror

import (
	"reflect"
	"testing"
)

func TestBasics(t *testing.T) {
	expectBasic := []interface{}{
		int8(0), uint8(0), int(0), int64(0),
		float64(0), complex128(0),
		' ', " ",
	}
	for _, v := range expectBasic {
		tv := reflect.TypeOf(v)
		Walk(tv, func(typ *reflect.StructField, typeIndex, depth int) error {
			if tv != typ.Type {
				t.Errorf("%v != %v\n", tv, typ.Type)
			}
			return nil
		})
	}
}
