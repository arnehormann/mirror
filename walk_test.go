package mirror

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unsafe" // only for unsafe.Pointer declaration  in test
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
	type embedded0 struct{}
	type embedded1 struct {
		a uint8
		b uint16
		c uint32
		d uint64
		e uint
		f func(interface {
			Func(uint)
		})
	}
	type embedded2 struct {
		a int8
		b int16
		c int32
		d int64
		e int `etag`
	}
	type compoundTest struct {
		embedded0
		embedded1
		a error `atag`
		_ [][][2]byte
		b *map[rune]*<-chan [2]uintptr
		c struct {
			a *complex64
			b complex128
			c interface{} `ctag`
			d interface {
				Do1()
				Do2() unsafe.Pointer
				Do3(func() error)
				Do4(<-chan [2]uintptr, chan<- [2]uintptr) (bool, float32, float64)
			}
			e func(string, int) (bool, uint16)
			f map[struct{}]interface{}
		}
		embedded2
		d struct{}
	}
	value := &compoundTest{}
	var actual string
	Walk(reflect.TypeOf(value), func(typ *reflect.StructField, typeIndex, depth int) error {
		actual += fmt.Sprintf("\n%d:%v,%v,%v,%v,%v,%v",
			depth, typ.Type.Kind(), typ.Type, typ.Index, typ.Name, typ.Tag, typ.Anonymous)
		return nil
	})
	// what we expect from the compoundTest types
	typeStrings := []string{
		`0:ptr,*mirror.compoundTest,[],,,false`,
		`1:struct,mirror.compoundTest,[],,,false`,
		`2:struct,mirror.embedded0,[0],embedded0,,true`,
		`2:struct,mirror.embedded1,[1],embedded1,,true`,
		`3:uint8,uint8,[0],a,,false`,
		`3:uint16,uint16,[1],b,,false`,
		`3:uint32,uint32,[2],c,,false`,
		`3:uint64,uint64,[3],d,,false`,
		`3:uint,uint,[4],e,,false`,
		`3:func,func(interface { Func(uint) }),[5],f,,false`,
		`4:interface,interface { Func(uint) },[],,,false`,
		`5:func,func(uint),[-1 0],Func,,false`,
		`6:uint,uint,[],,,false`,
		`2:interface,error,[2],a,atag,false`,
		`3:func,func() string,[-1 0],Error,,false`,
		`4:string,string,[],,,false`,
		`2:slice,[][][2]uint8,[3],_,,false`,
		`3:slice,[][2]uint8,[],,,false`,
		`4:array,[2]uint8,[],,,false`,
		`5:uint8,uint8,[],,,false`,
		`2:ptr,*map[int32]*<-chan [2]uintptr,[4],b,,false`,
		`3:map,map[int32]*<-chan [2]uintptr,[],,,false`,
		`4:int32,int32,[],,,false`,
		`4:ptr,*<-chan [2]uintptr,[],,,false`,
		`5:chan,<-chan [2]uintptr,[],,,false`,
		`6:array,[2]uintptr,[],,,false`,
		`7:uintptr,uintptr,[],,,false`,
		`2:struct,struct { a *complex64; b complex128; c interface {} "ctag"; d interface { Do1(); Do2() unsafe.Pointer; Do3(func() error); Do4(<-chan [2]uintptr, chan<- [2]uintptr) (bool, float32, float64) }; e func(string, int) (bool, uint16); f map[struct {}]interface {} },[5],c,,false`,
		`3:ptr,*complex64,[0],a,,false`,
		`4:complex64,complex64,[],,,false`,
		`3:complex128,complex128,[1],b,,false`,
		`3:interface,interface {},[2],c,ctag,false`,
		`3:interface,interface { Do1(); Do2() unsafe.Pointer; Do3(func() error); Do4(<-chan [2]uintptr, chan<- [2]uintptr) (bool, float32, float64) },[3],d,,false`,
		`4:func,func(),[-1 0],Do1,,false`,
		`4:func,func() unsafe.Pointer,[-1 1],Do2,,false`,
		`5:unsafe.Pointer,unsafe.Pointer,[],,,false`,
		`4:func,func(func() error),[-1 2],Do3,,false`,
		`5:func,func() error,[],,,false`,
		`6:interface,error,[],,,false`,
		`4:func,func(<-chan [2]uintptr, chan<- [2]uintptr) (bool, float32, float64),[-1 3],Do4,,false`,
		`5:chan,<-chan [2]uintptr,[],,,false`,
		`5:chan,chan<- [2]uintptr,[],,,false`,
		`6:array,[2]uintptr,[],,,false`,
		`5:bool,bool,[],,,false`,
		`5:float32,float32,[],,,false`,
		`5:float64,float64,[],,,false`,
		`3:func,func(string, int) (bool, uint16),[4],e,,false`,
		`4:string,string,[],,,false`,
		`4:int,int,[],,,false`,
		`4:bool,bool,[],,,false`,
		`4:uint16,uint16,[],,,false`,
		`3:map,map[struct {}]interface {},[5],f,,false`,
		`4:struct,struct {},[],,,false`,
		`4:interface,interface {},[],,,false`,
		`2:struct,mirror.embedded2,[6],embedded2,,true`,
		`3:int8,int8,[0],a,,false`,
		`3:int16,int16,[1],b,,false`,
		`3:int32,int32,[2],c,,false`,
		`3:int64,int64,[3],d,,false`,
		`3:int,int,[4],e,etag,false`,
		`2:struct,struct {},[7],d,,false`,
	}
	// used instead of a literal string to avoid conversion errors by git (\n to \r\n)
	expected := "\n" + strings.Join(typeStrings, "\n")
	if actual != expected {
		t.Error("walked compound type did not match")
	}
}
