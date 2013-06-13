package mirror

import (
	"reflect"
)

// VisitType processes type information.
// It is called by Walk.
// - typ is:
//   - a wrapped type (if len(typ.Index) == 0)
//   - an interface method (type.Index[0] < 0; ignore first value); this is clumsy...
//   - a field of a struct
// - typeIndex is the unique index of this type in the set types enountered during a walk
// - depth counts the number of indirections used to get to this type from the root type
//
// - for a map, the key and the element types are visited
// - for a function, the input argument types and output argument types are visited
// - for an interface, all contained functions are visited
type VisitType func(typ *reflect.StructField, typeIndex, depth int) error

type typeWalker map[reflect.Type]int

func (walker typeWalker) index(t reflect.Type) (int, bool) {
	idx, known := walker[t]
	if !known {
		idx = len(walker)
		walker[t] = idx
	}
	return idx, known
}

func (walker typeWalker) walk(t *reflect.StructField, visit VisitType) error {
	if t == nil || visit == nil {
		return nil
	}
	type stackNode struct {
		field      *reflect.StructField
		currentIdx int
		depth      int
		known      bool
	}
	idx, known := walker.index(t.Type)
	stack := []stackNode{{t, idx, 0, known}}
	var node stackNode
	for lastIdx := 0; lastIdx >= 0; lastIdx = len(stack) - 1 {
		stack, node = stack[:lastIdx], stack[lastIdx]
		err := visit(node.field, node.currentIdx, node.depth)
		if err != nil {
			return err
		}
		if node.known {
			// to prevent endless loops, don't follow known types
			continue
		}
		// follow container types
		depth := node.depth + 1
		switch t := node.field.Type; t.Kind() {
		case reflect.Struct:
			fields := t.NumField()
			for i := fields - 1; i >= 0; i-- {
				// add each field from struct to stack. As this is LIFO, start with the last field.
				field := t.Field(i)
				typeIdx, known := walker.index(field.Type)
				stack = append(stack, stackNode{&field, typeIdx, depth, known})
			}
		case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Chan, reflect.Map:
			// add element to stack
			field := &reflect.StructField{Type: t.Elem()}
			typeIdx, known := walker.index(field.Type)
			stack = append(stack, stackNode{field, typeIdx, depth, known})
			if t.Kind() == reflect.Map {
				// add key to stack
				field := &reflect.StructField{Type: t.Key()}
				typeIdx, known := walker.index(field.Type)
				stack = append(stack, stackNode{field, typeIdx, depth, known})
			}
		case reflect.Interface:
			for i := t.NumMethod() - 1; i >= 0; i-- {
				m := t.Method(i)
				typeIdx, known := walker.index(m.Type)
				stack = append(stack, stackNode{
					&reflect.StructField{
						Name:    m.Name,
						PkgPath: m.PkgPath,
						Type:    m.Type,
						Index:   []int{-1, m.Index},
					},
					typeIdx,
					depth,
					known,
				})
			}
		case reflect.Func:
			for i := t.NumOut() - 1; i >= 0; i-- {
				ret := t.Out(i)
				typeIdx, known := walker.index(ret)
				stack = append(stack, stackNode{
					&reflect.StructField{Type: ret},
					typeIdx,
					depth,
					known,
				})
			}
			for i := t.NumIn() - 1; i >= 0; i-- {
				arg := t.In(i)
				typeIdx, known := walker.index(arg)
				stack = append(stack, stackNode{
					&reflect.StructField{Type: arg},
					typeIdx,
					depth,
					known,
				})
			}
		}
	}
	return nil
}

// Walk will call visit on t and each type in t.
// It will follow struct fields, pointers, arrays, slices, maps and chans
// and call visit on each element.
// Walk will only follow into types at their very first occurence during this walk.
func Walk(t reflect.Type, visit VisitType) error {
	walker := make(typeWalker)
	return walker.walk(&reflect.StructField{Type: t}, visit)
}
