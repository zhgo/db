// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"reflect"
)

// Table struct
type Table struct {
	// Table name
	Name string

	// Table primary
	Primary string

	// All fields, except primary
	Fields []string

	// All fields
	AllFields []string

	// json and field property map
	FiledsMap map[string]string

	// Entity
	EntityType reflect.Type
}

// New Table
func NewTable(tableName string, entity interface{}) *Table {
	p, f, af, jf := tableFields(entity)
	return &Table{
		Name:       tableName,
		Primary:    p,
		Fields:     f,
		AllFields:  af,
		FiledsMap:  jf,
		EntityType: reflect.ValueOf(entity).Elem().Type(),
	}
}
