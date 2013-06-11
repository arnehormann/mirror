// mirror - helps to make unsafe a little safer
//
// Copyright 2013 Arne Hormann. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

// This file is generated to keep it internally consistent,
// do not edit it directly! You probably shouldn't change it.
// The generator is run from the branch 'taggenerator'
// with './create_tags.sh'

package mirror

import (
	"reflect"
	"strings"
)

// binary representation of tag based configuration
type tagflag uint

const (
	// flags and names

	tfIgnore tagflag = 1 << 0
	tnIgnore         = "ignore"

	tfFieldname tagflag = 1 << 1
	tnFieldname         = "fieldname"

	tfNofieldname tagflag = 1 << 2
	tnNofieldname         = "nofieldname"

	tfTypename tagflag = 1 << 3
	tnTypename         = "typename"

	tfNotypename tagflag = 1 << 4
	tnNotypename         = "notypename"

	tfType tagflag = 1 << 5
	tnType         = "type"

	tfNotype tagflag = 1 << 6
	tnNotype         = "notype"

	tfFollow tagflag = 1 << 7
	tnFollow         = "follow"

	tfNofollow tagflag = 1 << 8
	tnNofollow         = "nofollow"

	tfSame tagflag = 1 << 9
	tnSame         = "same"

	tfMethods tagflag = 1 << 10
	tnMethods         = "methods"

	tfAssignable tagflag = 1 << 11
	tnAssignable         = "assignable"

	// number of defined tags
	numtags = 12

	// masks
	tmIgnore      = ^tfIgnore
	tmFieldname   = ^tfNofieldname
	tmNofieldname = ^tfFieldname
	tmTypename    = ^tfNotypename
	tmNotypename  = ^tfTypename
	tmType        = ^tfNotype
	tmNotype      = ^tfType
	tmFollow      = ^tfNofollow
	tmNofollow    = ^tfFollow
	tmSame        = ^tfMethods & ^tfAssignable
	tmMethods     = ^tfSame & ^tfAssignable
	tmAssignable  = ^tfSame & ^tfMethods
)

// flagname converts a single flag to its string representation.
// keep in sync with tag-constants above, otherwise init panics
func (t tagflag) name() string {
	switch t {
	case tfIgnore:
		return tnIgnore
	case tfFieldname:
		return tnFieldname
	case tfNofieldname:
		return tnNofieldname
	case tfTypename:
		return tnTypename
	case tfNotypename:
		return tnNotypename
	case tfType:
		return tnType
	case tfNotype:
		return tnNotype
	case tfFollow:
		return tnFollow
	case tfNofollow:
		return tnNofollow
	case tfSame:
		return tnSame
	case tfMethods:
		return tnMethods
	case tfAssignable:
		return tnAssignable
	}
	// is not a single flag
	return ""
}

// toggle mutually exclusive values, intended for default flags
// does not do error checking
func (t tagflag) add(setTo tagflag) tagflag {
	// check each tag and mask forbidden others
	if setTo&tfIgnore == tfIgnore {
		t = t & ^tmIgnore
	}
	if setTo&tfFieldname == tfFieldname {
		t = t & ^tmFieldname
	}
	if setTo&tfNofieldname == tfNofieldname {
		t = t & ^tmNofieldname
	}
	if setTo&tfTypename == tfTypename {
		t = t & ^tmTypename
	}
	if setTo&tfNotypename == tfNotypename {
		t = t & ^tmNotypename
	}
	if setTo&tfType == tfType {
		t = t & ^tmType
	}
	if setTo&tfNotype == tfNotype {
		t = t & ^tmNotype
	}
	if setTo&tfFollow == tfFollow {
		t = t & ^tmFollow
	}
	if setTo&tfNofollow == tfNofollow {
		t = t & ^tmNofollow
	}
	if setTo&tfSame == tfSame {
		t = t & ^tmSame
	}
	if setTo&tfMethods == tfMethods {
		t = t & ^tmMethods
	}
	if setTo&tfAssignable == tfAssignable {
		t = t & ^tmAssignable
	}
	return t
}

// get all tagflags that are set though they must not be.
// No conflicts returns 0
func (t tagflag) conflicts() tagflag {
	return t & ^t.add(t)
}

// Has Tag "ignore" set
func (t tagflag) taggedIgnore() bool {
	return t&tfIgnore == tfIgnore
}

// Has Tag "fieldname" set
func (t tagflag) taggedFieldname() bool {
	return t&tfFieldname == tfFieldname
}

// Has Tag "nofieldname" set
func (t tagflag) taggedNofieldname() bool {
	return t&tfNofieldname == tfNofieldname
}

// Has Tag "typename" set
func (t tagflag) taggedTypename() bool {
	return t&tfTypename == tfTypename
}

// Has Tag "notypename" set
func (t tagflag) taggedNotypename() bool {
	return t&tfNotypename == tfNotypename
}

// Has Tag "type" set
func (t tagflag) taggedType() bool {
	return t&tfType == tfType
}

// Has Tag "notype" set
func (t tagflag) taggedNotype() bool {
	return t&tfNotype == tfNotype
}

// Has Tag "follow" set
func (t tagflag) taggedFollow() bool {
	return t&tfFollow == tfFollow
}

// Has Tag "nofollow" set
func (t tagflag) taggedNofollow() bool {
	return t&tfNofollow == tfNofollow
}

// Has Tag "same" set
func (t tagflag) taggedSame() bool {
	return t&tfSame == tfSame
}

// Has Tag "methods" set
func (t tagflag) taggedMethods() bool {
	return t&tfMethods == tfMethods
}

// Has Tag "assignable" set
func (t tagflag) taggedAssignable() bool {
	return t&tfAssignable == tfAssignable
}

// appliesTo returns true if all set tags are applicable to k.
func (t tagflag) appliesTo(k reflect.Kind) bool {
	switch {
	case t&tfFollow == tfFollow:
		return k == reflect.Struct
	case t&tfNofollow == tfNofollow:
		return k == reflect.Struct
	case t&tfSame == tfSame:
		return k == reflect.Interface
	case t&tfMethods == tfMethods:
		return k == reflect.Interface
	case t&tfAssignable == tfAssignable:
		return k == reflect.Interface
	default:
	}
	return true
}

func (t tagflag) String() string {
	s := ""
	for i := uint(0); i < numtags; i++ {
		if flag := tagflag(1 << i); t&flag != 0 {
			s += flag.name() + ","
		}
	}
	if s == "" {
		return ""
	}
	return s[:len(s)-1]
}

func (t tagflag) stringWithDefaults(defaults tagflag) string {
	return "tags[" + (t & ^defaults).String() + "] and defaults [" + (defaults & ^t).String() + "]"
}

// parseTag converts a string to its flag representation.
func parseTag(tag string) tagflag {
	switch tag {
	case tnIgnore:
		return tfIgnore
	case tnFieldname:
		return tfFieldname
	case tnNofieldname:
		return tfNofieldname
	case tnTypename:
		return tfTypename
	case tnNotypename:
		return tfNotypename
	case tnType:
		return tfType
	case tnNotype:
		return tfNotype
	case tnFollow:
		return tfFollow
	case tnNofollow:
		return tfNofollow
	case tnSame:
		return tfSame
	case tnMethods:
		return tfMethods
	case tnAssignable:
		return tfAssignable
	}
	return 0
}

type tagError string

func (e tagError) Error() string {
	return string(e)
}

func parse(seed tagflag, tags string) (tagflag, error) {
	// from StructField.Tag.Get("mirror")
	t := seed
	for _, tag := range strings.Split(",", tags) {
		tf := parseTag(tag)
		if tf == 0 {
			return 0, tagError("Unknown tag '" + tag + "'")
		}
		t |= tf
	}
	if c := t.conflicts(); c != 0 {
		return 0, tagError(t.stringWithDefaults(seed) + " conflict with [" + c.String() + "]")
	}
	return t, nil
}
