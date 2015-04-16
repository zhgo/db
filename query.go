// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
    "strings"
    "fmt"
    "errors"
    "log"
    "database/sql"
)

// Query type
const (
    QueryInsert = iota
    QueryUpdate
    QueryDelete
    QuerySelect
)

// Query struct
type Query struct {
    // Server
    Server *Server

    // Query type: Select, Insert, Update, Delete
    Type uint

    // Sql
    Sql map[string]string

    //args for Sql
    Args []interface{}

    // Limit
    limit QueryLimit

    // Current Sql node
    Current string
}

// Limit struct
type QueryLimit struct{
    // Enabled
    Enabled bool

    // Offset
    Offset int64

    // Rows
    Rows int64
}

// And
func (q *Query) And(f string, co string, v ...string) *Query {
    return q.Condition(q.Current, "AND", f, co, v...)
}

// Or
func (q *Query) Or(f string, co string, v ...string) *Query {
    return q.Condition(q.Current, "OR", f, co, v...)
}

// Not
func (q *Query) Not(f string, co string, v ...string) *Query {
    return q.Condition(q.Current, "NOT", f, co, v...)
}

// node Sql node
// l Logical
// f Field
// co Comparison Operators
// v Value(s)
func (q *Query) Condition(node string, l string, f string, co string, v ...string) *Query {
    q.Sql[node] += fmt.Sprintf(" %s %s %s '%s' ", l, f, co, strings.Join(v, "', '"))
    q.Current = node
    return q
}

// Select fields
func (q *Query) Select(f ...string) *Query {
    q.Type = QuerySelect
    if len(f) == 0 {
        q.Sql["Select"] = "*"
    } else {
        q.Sql["Select"] = strings.Join(f, ", ")
    }
    return q
}

// From
func (q *Query) From(tb string, alias ...string) *Query {
    q.Sql["From"] = fmt.Sprintf(" %s %s ", tb, tableAlias(alias))
    q.Current = "From"
    return q
}

// Join
func (q *Query) Join(tb string, alias ...string) *Query {
    return q.JoinNode("", tb, alias...)
}

// Join Inner
func (q *Query) JoinInner(tb string, alias ...string) *Query {
    return q.JoinNode("INNER", tb, alias...)
}

// Join Outer
func (q *Query) JoinOuter(tb string, alias ...string) *Query {
    return q.JoinNode("OUTER", tb, alias...)
}

// Join Left
func (q *Query) JoinLeft(tb string, alias ...string) *Query {
    return q.JoinNode("LEFT", tb, alias...)
}

// Join Right
func (q *Query) JoinRight(tb string, alias ...string) *Query {
    return q.JoinNode("RIGHT", tb, alias...)
}

// Join node
func (q *Query) JoinNode(p string, tb string, alias ...string) *Query {
    q.Sql["Join"] += fmt.Sprintf(" %s JOIN %s %s ", p, tb, tableAlias(alias))
    q.Current = "Join"
    return q
}

// Join On
func (q *Query) On(f string, co string, v ...string) *Query {
    return q.OnNode("ON", f, co, v...)
}

// On node
func (q *Query) OnNode(l string, f string, co string, v ...string) *Query {
    return q.Condition("Join", l, f, co, v...)
}

// Where
func (q *Query) Where(f string, co string, v ...string) *Query {
    return q.WhereNode("WHERE", f, co, v...)
}

// Where node
func (q *Query) WhereNode(l string, f string, co string, v ...string) *Query {
    return q.Condition("Where", l, f, co, v...)
}

// Group
func (q *Query) Group(f ...string) *Query {
    q.Sql["Group"] = fmt.Sprintf(" GROUP BY %s ", strings.Join(f, ", "))
    return q
}

// Having
func (q *Query) Having(f string, co string, v ...string) *Query {
    return q.Condition("Having", "HAVING", f, co, v...)
}

// Order
func (q *Query) Order(sort string, f ...string) *Query {
    q.Sql["Order"] = fmt.Sprintf(" ORDER BY %s %s ", strings.Join(f, ", "), sort)
    q.Current = "Order"
    return q
}

// Order ASC
func (q *Query) OrderAsc(f ...string) *Query {
    return q.Order("ASC", f...)
}

// Order DESC
func (q *Query) OrderDesc(f ...string) *Query {
    return q.Order("DESC", f...)
}

// Limit
func (q *Query) Limit(offset, rows int64) *Query {
    q.limit.Enabled = true
    q.limit.Offset = offset
    q.limit.Rows = rows
    q.Current = "Limit"
    return q
}

// Delete
func (q *Query) DeleteFrom(tb string) *Query {
    q.Type = QueryDelete
    q.Sql["Delete"] = tb
    q.Current = "Delete"
    return q
}

// Update
func (q *Query) Update(tb string) *Query {
    q.Type = QueryUpdate
    q.Sql["Update"] = tb
    q.Current = "Update"
    return q
}

// Set(Update)
func (q *Query) Set(f string, v interface{}) *Query {
    str, ok := q.Sql["UpdateSet"]
    if ok && len(str) > 0 {
        q.Sql["UpdateSet"] += fmt.Sprintf(" , %s = '%v' ", f, v)
    } else {
        q.Sql["UpdateSet"] = fmt.Sprintf(" SET %s = '%v' ", f, v)
    }
    q.Current = "UpdateSet"
    return q
}

// Insert
func (q *Query) InsertInto(tb string) *Query {
    q.Type = QueryInsert
    q.Sql["Insert"] = tb
    q.Current = "Insert"
    return q
}

// Fields(Insert)
func (q *Query) Fields(f ...string) *Query {
    q.Sql["InsertFields"] = fmt.Sprintf(" (%s) ", strings.Join(f, ", "))
    q.Current = "InsertFields"
    return q
}

// Values(Insert)
func (q *Query) Values(v ...string) *Query {
    str, ok := q.Sql["InsertValues"]
    if ok && len(str) > 0 {
        q.Sql["InsertValues"] += fmt.Sprintf(" ,('%v') ", strings.Join(v, "', '"))
    } else {
        q.Sql["InsertValues"] = fmt.Sprintf(" VALUES('%v') ", strings.Join(v, "', '"))
    }
    q.Current = "InsertValues"
    return q
}

// Connect all sql part to a corect sql string.
func (q *Query) toString() string {
    str := ""

    // Select
    if node, ok := q.Sql["Select"]; ok && len(node) > 0 {
        str += fmt.Sprintf(" SELECT %s ", node)
    }
    if node, ok := q.Sql["From"]; ok && len(node) > 0 {
        str += fmt.Sprintf(" From %s ", node)
    }
    if node, ok := q.Sql["Join"]; ok && len(node) > 0 {
        str += fmt.Sprintf(" %s ", node)
    }

    // Insert
    if node, ok := q.Sql["Insert"]; ok && len(node) > 0 {
        str += fmt.Sprintf(" INSERT INTO %s ", node)
    }
    if node, ok := q.Sql["InsertFields"]; ok && len(node) > 0 {
        str += fmt.Sprintf(" %s ", node)
    }
    if node, ok := q.Sql["InsertValues"]; ok && len(node) > 0 {
        str += fmt.Sprintf(" %s ", node)
    }

    // Update
    if node, ok := q.Sql["Update"]; ok && len(node) > 0 {
        str += fmt.Sprintf(" UPDATE %s ", node)
    }
    if node, ok := q.Sql["UpdateSet"]; ok && len(node) > 0 {
        str += fmt.Sprintf(" %s ", node)
    }

    // Delete
    if node, ok := q.Sql["Delete"]; ok && len(node) > 0 {
        str += fmt.Sprintf(" DELETE FROM %s ", node)
    }

    // Select, Update, Delete
    if node, ok := q.Sql["Where"]; ok && len(node) > 0 {
        str += fmt.Sprintf(" %s ", node)
    }

    // Select
    if node, ok := q.Sql["Group"]; ok && len(node) > 0 {
        str += fmt.Sprintf(" %s ", node)
    }
    if node, ok := q.Sql["Having"]; ok && len(node) > 0 {
        str += fmt.Sprintf(" %s ", node)
    }

    // Select, Update
    if node, ok := q.Sql["Order"]; ok && len(node) > 0 {
        str += fmt.Sprintf(" %s ", node)
    }

    // Select, Update
    if q.limit.Enabled == true {
        str += fmt.Sprintf(" LIMIT %d, %d ", q.limit.Offset, q.limit.Rows)
    }

    // Select
    if node, ok := q.Sql["ForUpdate"]; ok && len(node) > 0 {
        str += fmt.Sprintf(" %s ", node)
    }

    log.Printf("%s\n", str)
    log.Printf("%#v\n", q.Args)

    return str
}

// Exec
func (q *Query) Exec() (sql.Result, error) {
    if q.Server == nil {
        return nil, errors.New("DB config not found")
    }

    return q.Server.Exec(q.toString(), q.Args...)
}

// Row
func (q *Query) Row(ptr interface{}) error {
    if q.Server == nil {
        return errors.New("DB config not found")
    }

    return q.Server.Row(ptr, q.toString(), q.Args...)
}

// Rows
func (q *Query) Rows(ptr interface{}) error {
    if q.Server == nil {
        return errors.New("DB config not found")
    }

    return q.Server.Rows(ptr, q.toString(), q.Args...)
}

// Drop Table
func (q *Query) DropTable(tb string) error {
    if q.Server == nil {
        return errors.New("DB config not found")
    }

    _, err := q.Server.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tb))
    return err
}

// New Query object
func NewQuery(server *Server) *Query {
    query := Query{Server: server}
    query.Sql = make(map[string]string)
    query.Args = make([]interface{}, 0)
    return &query
}
