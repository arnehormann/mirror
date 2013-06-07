// mirror - helps to make unsafe a little safer
//
// Copyright 2013 Arne Hormann. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

// +build ignore

// generate tags with 'create_tags.sh'

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

var _ = fmt.Errorf // import fmt for string debugging even though we don't use it all the time

var allTagGroups = []string{
	// parseable tags, details may be configured by postfix
	// the postfix meanings are:
	// - "_" default for variable name "_"
	// - "*" default for named variables
	// - "^" exclusive, forbid all other flags
	// In addition, an element in "[" ... "]" sets the reflect.Kind for the group
	// implemented in expandTagGroups
	"ignore^",
	"fieldname* nofieldname_",
	"typename* notypename_",
	"type* notype_",
	"[Struct] follow* nofollow",
	"[Interface] same methods assignable*",
}

type tag struct {
	Name      string
	Forbidden []string
	// for default constant generation
	Kind         string
	DefaultNamed bool
	DefaultBlank bool
	// for tags prohibiting all others (just "ignore" now)
	Exclusive bool
}

func prepareTag(rawtag string) (t *tag) {
	t = &tag{}
	name := rawtag
	for len(name) > 0 {
		lastIdx := len(name) - 1
		switch postfix := name[lastIdx:]; postfix {
		case "*":
			t.DefaultNamed = true
		case "_":
			t.DefaultBlank = true
		case "^":
			t.Exclusive = true
		case "]":
			t.Kind = name[1:lastIdx]
			return
		default:
			t.Name = name
			return
		}
		name = name[:lastIdx]
	}
	return
}

func expandTagGroups(tagGroups []string) []tag {
	knownTags := make(map[string]*tag)
	tagnames := make([]string, 0, 16)
	for _, tagGroup := range tagGroups {
		var groupMembers []*tag
		var kind string
		groupTagnames := strings.Split(tagGroup, " ")
		for _, tagDesc := range groupTagnames {
			newTag := prepareTag(tagDesc)
			switch {
			case newTag.Kind != "":
				kind = newTag.Kind
			case newTag.Name != "":
				groupMembers = append(groupMembers, newTag)
			}
		}
		for i, currentTag := range groupMembers {
			forbidden := make([]string, 0, len(groupMembers)-1)
			for j, other := range groupMembers {
				if i != j {
					forbidden = append(forbidden, other.Name)
				}
			}
			currentTag.Forbidden = forbidden
			currentTag.Kind = kind
			knownTags[currentTag.Name] = currentTag
			tagnames = append(tagnames, currentTag.Name)
		}
	}
	tags := make([]tag, len(tagnames))
	for i, name := range tagnames {
		currentTag := knownTags[name]
		if currentTag.Exclusive {
			currentTag.Forbidden = append(append([]string{}, tagnames[:i]...), tagnames[i+1:]...)
		}
		tags[i] = *currentTag
	}
	return tags
}

func main() {
	// Template data
	tags := expandTagGroups(allTagGroups)
	// Template
	funcMap := template.FuncMap{
		"title": strings.Title,
		"flag": func(v string) string {
			return "tf" + strings.Title(v)
		},
		"name": func(v string) string {
			return "tn" + strings.Title(v)
		},
		"mask": func(v string) string {
			return "tm" + strings.Title(v)
		},
	}
	tagTemplate, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	templater := template.New("tags.go").Delims("~", "~").Funcs(funcMap)
	templater = template.Must(templater.Parse(string(tagTemplate)))
	if err := templater.Execute(os.Stdout, tags); err != nil {
		panic(err)
	}
}
