// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"errors"
	"log"
	"reflect"
)

// Item
type Item map[string]interface{}

// Where
type Where map[string]interface{}

var Env int8 = 0

// Get scan variables
func scanVariables(ptr interface{}, columnsLen int, isRows bool) (reflect.Kind, interface{}, []interface{}, error) {
	typ := reflect.ValueOf(ptr).Type()

	if typ.Kind() != reflect.Ptr {
		return 0, nil, nil, errors.New("ptr is not a pointer")
	}

	//log.Printf("%s\n", dataType.Elem().Kind())
	elemTyp := typ.Elem()

	if isRows { // Rows
		if elemTyp.Kind() != reflect.Slice {
			return 0, nil, nil, errors.New("ptr is not point a slice")
		}

		elemTyp = elemTyp.Elem()
	}

	elemKind := elemTyp.Kind()

	// element(value) is point to row
	scan := make([]interface{}, columnsLen)

	//log.Printf("%s\n", elemKind)

	if elemKind == reflect.Struct {
		if columnsLen != elemTyp.NumField() {
			return 0, nil, nil, errors.New("columnsLen is not equal elemTyp.NumField()")
		}

		row := reflect.New(elemTyp) // Data
		for i := 0; i < columnsLen; i++ {
			f := elemTyp.Field(i)
			if !f.Anonymous { // && f.Tag.Get("json") != ""
				scan[i] = row.Elem().FieldByIndex([]int{i}).Addr().Interface()
			}
		}

		return elemKind, row.Interface(), scan, nil
	}

	if elemKind == reflect.Map || elemKind == reflect.Slice {
		row := make([]interface{}, columnsLen) // Data
		for i := 0; i < columnsLen; i++ {
			scan[i] = &row[i]
		}

		return elemKind, &row, scan, nil
	}

	return 0, nil, nil, errors.New("ptr is not a point struct, map or slice")
}

// Type assertions
func typeAssertion(v interface{}) interface{} {
	switch v.(type) {
	case []byte:
		return v.([]byte)
	case []rune:
		return v.([]rune)
	case bool:
		return v.(bool)
	case float64:
		return v.(float64)
	case int64:
		return v.(int64)
	case nil:
		return nil
	case string:
		return v.(string)
	default:
		log.Printf("Unexpected type %#v\n", v)
		return ""
	}
}

// Table alias
func tableAlias(alias []string) string {
	if len(alias) > 0 {
		return alias[0]
	}
	return ""
}

// Reflect struct, construct Field slice
func tableFields(entity interface{}) (string, []string) {
	typ := reflect.Indirect(reflect.ValueOf(entity)).Type()
	primary := ""
	fields := make([]string, 0)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		var name string
		if field.Tag.Get("json") != "" {
			name = field.Tag.Get("json")
		} else {
			name = field.Name
		}

		if field.Tag.Get("pk") == "true" { //!field.Anonymous
			primary = name
		} else {
			fields = append(fields, name)
		}
	}

	return primary, fields
}
