// Copyright 2014 The fav Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"errors"
	"log"
	"reflect"
)

//DB config
type Config struct {
	//Database type
	Type string

	//DSN
	DSN string

	//profiling
	Profiling bool
}

//sql struct
type Sql struct {
	//Select fields
	Select string

	//FROM Cation: no use normal
	From string

	//Table join
	Join string

	//Where
	Where string

	//Order by
	Order string

	//Group by
	Group string

	//having
	Having string

	//limit x
	Offset uint64 //limit Offset, Rows

	//limit x, y
	Rows uint64 //limit Offset, Rows

	//lock row(s)
	ForUpdate string

	//args for sql
	Args []interface{}
}

//Table struct
type Table struct {
	//Table name
	Name string

	//Table primary
	Primary string

	//All fields, except primary
	Fields []string
}

//New query
func NewQuery(table Table, dbKey string, config Config) *Query {
	return &Query{table: table, sql: Sql{Offset: 0, Rows: 10}, execute: Execute{DBKey: dbKey, Config: config}}
}

//分析传入指针变量,构建rows.Scan所必须的参数
//reflect.Kind 返回dataPtr类型
//interface{} 返回struct/alice指针
//interface{} 返回rows.Scan所须的指针切片
//error 返回错误
func scanVariable(dataPtr interface{}, columnsLen int, isRows bool) (reflect.Kind, interface{}, []interface{}, error) {
	dataValue := reflect.ValueOf(dataPtr)
	dataType := dataValue.Type()

	if dataType.Kind() != reflect.Ptr {
		return 0, nil, nil, errors.New("dataPtr is not a pointer")
	}

	//log.Printf("%s\n", dataType.Elem().Kind())

	if isRows && dataType.Elem().Kind() != reflect.Slice {
		return 0, nil, nil, errors.New("dataPtr is not point a slice")
	}

	//columnsLen := len(columns)

	scanArgs := make([]interface{}, columnsLen) //指针

	//var elemVal reflect.Value
	var elemTyp reflect.Type

	if isRows { //多行数据的情况
		//elemVal := dataValue.Elem().Elem()
		elemTyp = dataType.Elem().Elem()
	} else { //单行数据的情况
		//elemVal := dataValue.Elem()
		elemTyp = dataType.Elem()
	}

	elemKind := elemTyp.Kind()

	//log.Printf("%s\n", elemKind)

	if elemKind == reflect.Struct {
		elemNumField := elemTyp.NumField()
		if columnsLen != elemNumField {
			return 0, nil, nil, errors.New("columnsLen is not equal elemNumField")
		}

		scanVals := reflect.New(elemTyp)
		for i := 0; i < elemNumField; i++ {
			elemField := elemTyp.Field(i)
			if !elemField.Anonymous { // && elemField.Tag.Get("json") != ""
				scanArgs[i] = scanVals.Elem().FieldByIndex([]int{i}).Addr().Interface()
			}
		}
		return reflect.Struct, scanVals.Interface(), scanArgs, nil
	} else if elemKind == reflect.Map || elemKind == reflect.Slice {
		scanVals := make([]interface{}, columnsLen) //数据
		for i := 0; i < columnsLen; i++ {
			scanArgs[i] = &scanVals[i]
		}
		return elemKind, &scanVals, scanArgs, nil
	} else {
		return 0, nil, nil, errors.New("dataPtr is not point struct, map or slice")
	}
}

//Type assertions
func typeAssertion(v interface{}) interface{} {
	var r interface{}

	switch v.(type) {
	case bool:
		//log.Printf("bool\n")
		r = v.(bool)
	case int64:
		//log.Printf("int64\n")
		r = v.(int64)
	case float64:
		//log.Printf("float64\n")
		r = v.(float64)
	case string:
		//log.Printf("string\n")
		r = v.(string)
	case []byte:
		//log.Printf("[]byte\n")
		r = string(v.([]byte))
	default:
		log.Printf("Unexpected type %#v\n", v)
		r = ""
	}

	return r
}

//reflect struct, construct Field slice
func TableFields(entity interface{}) (string, []string) {
	typ := reflect.ValueOf(entity).Elem().Type()
	n := typ.NumField()

	primary := ""
	fields := make([]string, 0)

	for i := 0; i < n; i++ {
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
