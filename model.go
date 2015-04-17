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
func (m *Model) Insert() *Model {
    return &Model{Module: m.Module, Table: m.Table}
}

// Update
func (m *Model) Update() *Model {
    return &Model{Module: m.Module, Table: m.Table}
}

// Delete
func (m *Model) Delete() *Model {
    return &Model{Module: m.Module, Table: m.Table}
}

// Select
func (m *Model) Select() *Model {
    return &Model{Module: m.Module, Table: m.Table}
}

