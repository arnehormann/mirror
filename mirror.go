// mirror - helps to make unsafe a little safer
//
// Copyright 2013 Arne Hormann. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package mirror

import (
	"reflect"
)

// CanConvert returns true if the memory layout and the struct field names of
// 'from' match those of 'to'.
func CanConvert(from, to reflect.Type) bool {
	if from.Kind() != reflect.Struct || from.Kind() != to.Kind() ||
		from.Name() != to.Name() || from.NumField() != to.NumField() {
		return false
	}
	for i, max := 0, from.NumField(); i < max; i++ {
		sf, tf := from.Field(i), to.Field(i)
		if sf.Name != tf.Name || sf.Offset != tf.Offset {
			return false
		}
		tsf, ttf := sf.Type, tf.Type
		for done := false; !done; {
			k := tsf.Kind()
			if k != ttf.Kind() {
				return false
			}
			switch k {
			case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
				tsf, ttf = tsf.Elem(), ttf.Elem()
			case reflect.Interface:
				// don't have to handle matching interfaces here
				if tsf != ttf {
					// there are none in our case, so we are extra strict
					return false
				}
			case reflect.Struct:
				if recurseStructs <= 0 && tsf.Name() != ttf.Name() {
					return false
				}
				done = true
			default:
				done = true
			}
		}
		if recurseStructs > 0 && !CanConvert(tsf, ttf, recurseStructs-1) {
			return false
		}
	}
	return true
}
