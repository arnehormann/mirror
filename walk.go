package mirror

import (
	"reflect"
)

// VisitType
type VisitType func(t *reflect.StructField, typeIndex, depth int) error

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
	addNodes := func(stack []stackNode, nodes ...stackNode) []stackNode {
		result := append(stack, nodes...)
		copy(result[:len(stack)], result[len(nodes):])
		copy(nodes, stack[:len(nodes)])
		return result
	}
	idx, known := walker.index(t.Type)
	stack := []stackNode{{t, idx, idx, 0, known}}
	var node stackNode
	for len(stack) > 0 {
		node, stack = stack[0], stack[1:]
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
			for i, iLimit := 0, t.NumField(); i < iLimit; i++ {
				// add each field from struct to stack
				field := t.Field(i)
				typeIdx, known := walker.index(field.Type)
				stack = addNodes(stack, stackNode{&field, typeIdx, parentIdx, depth, known})
			}
		case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
			// add element to stack
			field := &reflect.StructField{Type: t.Elem()}
			typeIdx, known := walker.index(field.Type)
			stack = addNodes(stack, stackNode{field, typeIdx, parentIdx, depth, known})
		}
	}
	return nil
}

// Walk will call visit on t and each type in t.
// It will follow struct fields, pointers, arrays, slices, maps and chans
// and call visit on each element.
// For struct fields, Walk will only follow into fields at the first occurence of their type.
// Function handling for reflect.Interface and reflect.Method has to be done by visit.
func Walk(t reflect.Type, visit VisitType) error {
	walker := &typeWalker{indices: make(map[reflect.Type]int)}
	return walker.walk(&reflect.StructField{Type: t}, visit)
}