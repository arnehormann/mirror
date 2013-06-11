package mirror

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {
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

func TestCompound(t *testing.T) {
	type compoundTest struct {
		err   error `mytag`
		_     [][][1]byte
		c     *map[string]*<-chan [2]int
		inner struct {
			i *int
			v interface{} `tag`
		}
	}
	value := &compoundTest{}
	var asString string
	Walk(reflect.TypeOf(value), func(typ *reflect.StructField, typeIndex, depth int) error {
		asString += fmt.Sprintf("\n%d:%v,%v,%v,%v,%v,%v",
			depth, typ.Type.Kind(), typ.Type, typ.Index, typ.Name, typ.Tag, typ.Anonymous)
		return nil
	})
	typeStrings := []string{
		`0:ptr,*mirror.compoundTest,[],,,false`,
		`1:struct,mirror.compoundTest,[],,,false`,
		`2:interface,error,[0],err,mytag,false`,
		`2:slice,[][][1]uint8,[1],_,,false`,
		`3:slice,[][1]uint8,[],,,false`,
		`4:array,[1]uint8,[],,,false`,
		`5:uint8,uint8,[],,,false`,
		`2:ptr,*map[string]*<-chan [2]int,[2],c,,false`,
		`3:map,map[string]*<-chan [2]int,[],,,false`,
		`4:ptr,*<-chan [2]int,[],,,false`,
		`5:chan,<-chan [2]int,[],,,false`,
		`6:array,[2]int,[],,,false`,
		`7:int,int,[],,,false`,
		`2:struct,struct { i *int; v interface {} "tag" },[3],inner,,false`,
		`3:ptr,*int,[0],i,,false`,
		`4:int,int,[],,,false`,
		`3:interface,interface {},[1],v,tag,false`,
	}
	// we have to do it this way to avoid conversion errors by git (\n to \r\n)
	expect := "\n" + strings.Join(typeStrings, "\n")
	if asString != expect {
		t.Error("walked compound type did not match")
	}
}
