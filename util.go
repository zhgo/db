// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
    "errors"
    "log"
    "reflect"
)

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
        case bool:
        //log.Printf("bool\n")
        return v.(bool)
        case int64:
        //log.Printf("int64\n")
        return v.(int64)
        case float64:
        //log.Printf("float64\n")
        return v.(float64)
        case string:
        //log.Printf("string\n")
        return v.(string)
        case []byte:
        //log.Printf("[]byte\n")
        return string(v.([]byte))
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
