package mirror

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

/*
func TestDatabase(t *testing.T) {
	// not a real test yet...
	Walk(reflect.TypeOf(&sql.Row{}), func(t *reflect.StructField, typeIndex, depth int) error {
		fmt.Printf("%4d/%4d:\t%v\n", typeIndex, depth, t)
		return nil
	})
}
*/
