// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
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

// New Model
func NewModel(module string, table Table) *Model {
    return &Model{Module: module, Table: table}
}

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

