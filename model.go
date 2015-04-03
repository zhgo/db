// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
    "strings"
)

// Model struct
type Model struct {
    // Module name, as DB name.
    Module string

    // table instance
    Table Table
}

func NewModelQuery(m *Model) *Query {
    db, s := Databases[m.Module]
    if s == false {
        return nil //errors.New("DB config not found.")
    }

    return &Query{DB: db, Table: m.Table}
}

// Select row(s).
func (m *Model) Select(fields ...string) *Query {
    q := NewModelQuery(m)
    if len(fields) == 0 {
        fields = append(fields, "*")
    }

    q.sql.Select = strings.Join(fields, ", ")
    return q
}

// Insert.
func (m *Model) Insert(data *[]interface{}) (int64, error) {
    q := NewModelQuery(m)
    return q.Insert(data)
}

func (m *Model) Update(data *map[string]interface{}) (int64, error) {
    q := NewModelQuery(m)
    return q.Update(data)
}

func (m *Model) Delete() (int64, error) {
    q := NewModelQuery(m)
    return q.Delete()
}