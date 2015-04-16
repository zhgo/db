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

    // Entity
    EntityType reflect.Type
}

func NewTable(tableName string, entity interface{}) Table {
    p, f := tableFields(entity)
    t := Table{Name: tableName, Primary: p, Fields: f, EntityType: reflect.ValueOf(entity).Elem().Type()}
    return t
}