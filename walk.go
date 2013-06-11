package mirror

import (
	"reflect"
)

// VisitType processes type information.
// It is called by Walk.
// - typ is a field of a struct or a wrapped type (if len(typ.Index) == 0)
// - typeIndex is the unique index of this type in the set types enountered during a walk
// - depth counts the number of indirections used to get to this type from the root type
type VisitType func(typ *reflect.StructField, typeIndex, depth int) error

type typeWalker struct {
	indices map[reflect.Type]int
	types   []reflect.Type
}

func (walker *typeWalker) index(t reflect.Type) (int, bool) {
	idx, known := walker.indices[t]
	if known {
		return idx, true
	}
	idx = len(walker.types)
	walker.indices[t] = idx
	walker.types = append(walker.types, t)
	return idx, false
}

func (walker *typeWalker) walk(t *reflect.StructField, visit VisitType) error {
	if t == nil || visit == nil {
		return nil
	}
	type stackNode struct {
		field      *reflect.StructField
		currentIdx int
		parentIdx  int
		depth      int
		known      bool
	}
	idx, known := walker.index(t.Type)
	stack := []stackNode{{t, idx, idx, 0, known}}
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
		parentIdx := node.currentIdx
		depth := node.depth + 1
		switch t := node.field.Type; t.Kind() {
		case reflect.Struct:
			fields := t.NumField()
			for i := fields - 1; i >= 0; i-- {
				// add each field from struct to stack. As this is LIFO, start with the last field.
				field := t.Field(i)
				typeIdx, known := walker.index(field.Type)
				stack = append(stack, stackNode{&field, typeIdx, parentIdx, depth, known})
			}
		case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
			// add element to stack
			field := &reflect.StructField{Type: t.Elem()}
			typeIdx, known := walker.index(field.Type)
			stack = append(stack, stackNode{field, typeIdx, parentIdx, depth, known})
		}
	}
	return nil
}

// Walk will call visit on t and each type in t.
// It will follow struct fields, pointers, arrays, slices, maps and chans
// and call visit on each element.
// Walk will only follow into types at their very first occurence during this walk.
// Function handling for reflect.Interface and reflect.Method has to be done by visit.
// Map key types and function argument and return types are not processed.
func Walk(t reflect.Type, visit VisitType) error {
	walker := &typeWalker{indices: make(map[reflect.Type]int)}
	return walker.walk(&reflect.StructField{Type: t}, visit)
}
