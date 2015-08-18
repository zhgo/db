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

	// Fields for select, include primary
	SelectFields []string

	// Fields for add
	AddFields []string

	// Fields for update
	UpdateFields []string

	// json and field property map
	FiledsMap map[string]string

	// Entity type
	EntityType reflect.Type
}

// New Table
func NewTable(tableName string, entity interface{}) *Table {
	primary := ""
	fields := make([]string, 0)
	selectFields := make([]string, 0)
	addFields := make([]string, 0)
	updateFields := make([]string, 0)
	filedsMap := make(map[string]string)
	typ := reflect.Indirect(reflect.ValueOf(entity)).Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fd := field.Name
		if field.Tag.Get("field") != "" {
			fd = field.Tag.Get("field")
		}

		jn := field.Name
		if field.Tag.Get("json") != "" {
			jn = field.Tag.Get("json")
		}

		//!field.Anonymous
		if field.Tag.Get("pk") == "true" {
			primary = fd
		} else {
			fields = append(fields, fd)
		}

		selectFields = append(selectFields, fd)
		filedsMap[jn] = fd
	}

	return &Table{
		Name:         tableName,
		Primary:      primary,
		Fields:       fields,
		SelectFields: selectFields,
		AddFields:    addFields,
		UpdateFields: updateFields,
		FiledsMap:    filedsMap,
		EntityType:   typ,
	}
}
