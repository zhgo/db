// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
    "reflect"
)

// Model struct
type Model struct {
    // Module name, as DB name.
    Module string

    // table instance
    Table Table
}

// Server list
var Servers = make(map[string]*Server)

// Insert
func (m *Model) Insert() *Query {
    q := NewQuery(Servers[m.Module])
    q.InsertInto(m.Table.Name)
    return q
}

// Update
func (m *Model) Update() *Query {
    q := NewQuery(Servers[m.Module])
    q.Update(m.Table.Name)
    return q
}

// Delete
func (m *Model) Delete() *Query {
    q := NewQuery(Servers[m.Module])
    q.DeleteFrom(m.Table.Name)
    return q
}

// Select
func (m *Model) Select() *Query {
    q := NewQuery(Servers[m.Module])
    q.Select(m.Table.Fields...)
    q.From(m.Table.Name)
    return q
}

// New Model
func NewModel(module string, table Table) *Model {
    return &Model{Module: module, Table: table}
}


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

// New Table
func NewTable(tableName string, entity interface{}) Table {
    p, f := tableFields(entity)
    t := Table{Name: tableName, Primary: p, Fields: f, EntityType: reflect.ValueOf(entity).Elem().Type()}
    return t
}
