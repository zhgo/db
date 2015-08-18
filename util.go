// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"errors"
	"log"
	"reflect"
)

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

// Get scan variables
func scanVariables(ptr interface{}, columnsLen int, isRows bool) (reflect.Kind, interface{}, []interface{}, error) {
	typ := reflect.TypeOf(ptr)
	if typ.Kind() != reflect.Ptr {
		return 0, nil, nil, errors.New("ptr is not a pointer")
	}

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

	if elemKind == reflect.Struct {
		if columnsLen != elemTyp.NumField() {
			return 0, nil, nil, errors.New("columnsLen is not equal elemTyp.NumField()")
		}

		row := reflect.New(elemTyp) // Data
		for i := 0; i < columnsLen; i++ {
			f := elemTyp.Field(i)
			if !f.Anonymous { // && f.Tag.Get("field") != ""
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
