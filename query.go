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

// Query Nodes
var queryNodes = map[uint][]string{
    QueryInsert: []string{"Insert", "Fields", "Values"},
    QueryUpdate: []string{"Update", "Set", "Where", "Order", "Limit"},
    QueryDelete: []string{"Delete", "Where"},
    QuerySelect: []string{"Select", "From", "Join", "Where", "Group", "Having", "Order", "Limit", "ForUpdate"}}

// Query struct
type Query struct {
    // Server
    Server *Server

    // Query type: Insert, Update, Delete, Select
    Type uint

    // Sql
    Sql map[string]string

    //args for Sql
    Args []interface{}

    // Condition SQL
    sqlCond string

    // Current Sql node
    current string
}



// Equal
func (q *Query) Eq(f string, v interface{}) *Query {
    return q.condition("", f, "=", v)
}

// Greater than or equal
func (q *Query) Ge(f string, v interface{}) *Query {
    return q.condition("", f, ">=", v)
}

// Greater than
func (q *Query) Gt(f string, v interface{}) *Query {
    return q.condition("", f, ">", v)
}

// Less than or equal
func (q *Query) Le(f string, v interface{}) *Query {
    return q.condition("", f, "<=", v)
}

// Less than
func (q *Query) Lt(f string, v interface{}) *Query {
    return q.condition("", f, "<", v)
}

// Not equal
func (q *Query) Ne(f string, v interface{}) *Query {
    return q.condition("", f, "<>", v)
}

// Like: Simple pattern matching
func (q *Query) Like(f string, v interface{}) *Query {
    return q.condition("", f, "LIKE", v)
}

// In: Check whether a value is within a set of values
func (q *Query) In(f string, v ...interface{}) *Query {
    return q.conditionIn("", f, v...)
}



// and
func (q *Query) and(f string, co string, v interface{}) *Query {
    return q.condition("AND", f, co, v)
}

// And (Equal)
func (q *Query) AndEq(f string, v interface{}) *Query {
    return q.and(f, "=", v)
}

// And (Greater than or equal)
func (q *Query) AndGe(f string, v interface{}) *Query {
    return q.and(f, ">=", v)
}

// And (Greater than)
func (q *Query) AndGt(f string, v interface{}) *Query {
    return q.and(f, ">", v)
}

// And (Less than or equal)
func (q *Query) AndLe(f string, v interface{}) *Query {
    return q.and(f, "<=", v)
}

// And (Less than)
func (q *Query) AndLt(f string, v interface{}) *Query {
    return q.and(f, "<", v)
}

// And (Not equal)
func (q *Query) AndNe(f string, v interface{}) *Query {
    return q.and(f, "<>", v)
}

// And (Like: Simple pattern matching)
func (q *Query) AndLike(f string, v interface{}) *Query {
    return q.and(f, "LIKE", v)
}

// And (In: Check whether a value is within a set of values)
func (q *Query) AndIn(f string, v ...interface{}) *Query {
    return q.conditionIn("AND", f, v...)
}



// or
func (q *Query) or(f string, co string, v interface{}) *Query {
    return q.condition("OR", f, co, v)
}

// Or (Equal)
func (q *Query) OrEq(f string, v interface{}) *Query {
    return q.or(f, "=", v)
}

// Or (Greater than or equal)
func (q *Query) OrGe(f string, v interface{}) *Query {
    return q.or(f, ">=", v)
}

// Or (Greater than)
func (q *Query) OrGt(f string, v interface{}) *Query {
    return q.or(f, ">", v)
}

// Or (Less than or equal)
func (q *Query) OrLe(f string, v interface{}) *Query {
    return q.or(f, "<=", v)
}

// Or (Less than)
func (q *Query) OrLt(f string, v interface{}) *Query {
    return q.or(f, "<", v)
}

// Or (Not equal)
func (q *Query) OrNe(f string, v interface{}) *Query {
    return q.or(f, "<>", v)
}

// Or (Like: Simple pattern matching)
func (q *Query) OrLike(f string, v interface{}) *Query {
    return q.or(f, "LIKE", v)
}

// Or (In: Check whether a value is within a set of values)
func (q *Query) OrIn(f string, v ...interface{}) *Query {
    return q.conditionIn("OR", f, v...)
}


// l Logical
// f Field
// co Comparison Operators
// v Value
func (q *Query) condition(l string, f string, co string, v interface{}) *Query {
    q.sqlCond += fmt.Sprintf(" %s %s %s ? ", l, f, co)
    q.Args = append(q.Args, v)
    return q
}

// l Logical
// f Field
// vs Values
func (q *Query) conditionIn(l string, f string, vs ...interface{}) *Query {
    // Placeholder
    ph := make([]string, len(vs))
    for i, v := range vs {
        ph[i] = "?"
        q.Args = append(q.Args, v)
    }
    q.sqlCond += fmt.Sprintf(" %s %s IN (%s) ", l, f, strings.Join(ph, ", "))
    return q
}




// Select fields
func (q *Query) Select(f ...string) *Query {
    q.Type = QuerySelect
    if len(f) == 0 {
        q.Sql["Select"] = " SELECT *"
    } else {
        q.Sql["Select"] = fmt.Sprintf(" SELECT %s ", strings.Join(f, ", "))
    }
    q.current = "Select"
    return q
}

// From
func (q *Query) From(tb string) *Query {
    q.Sql["From"] = fmt.Sprintf(" FROM %s ", tb)
    q.current = "From"
    return q
}

// join
func (q *Query) join(p string, tb string) *Query {
    q.Sql["Join"] += fmt.Sprintf(" %s JOIN %s ", p, tb)
    q.current = "Join"
    return q
}

// Join
func (q *Query) Join(tb string) *Query {
    return q.join("", tb)
}

// Join Inner
func (q *Query) InnerJoin(tb string) *Query {
    return q.join("INNER", tb)
}

// Join Outer
func (q *Query) OuterJoin(tb string) *Query {
    return q.join("OUTER", tb)
}

// Join Left
func (q *Query) LeftJoin(tb string) *Query {
    return q.join("LEFT", tb)
}

// Join Right
func (q *Query) RightJoin(tb string) *Query {
    return q.join("RIGHT", tb)
}

// Join On
func (q *Query) On(qs ...*Query) *Query {
    q.Sql["Join"] += fmt.Sprintf(" ON %s ", q.sqlCond)
    q.sqlCond = ""
    q.current = "Join"
    return q
}



// Where
func (q *Query) Where(qs ...*Query) *Query {
    q.Sql["Where"] = fmt.Sprintf(" WHERE %s ", q.sqlCond)
    q.sqlCond = ""
    q.current = "Where"
    return q
}

// Group By
func (q *Query) GroupBy(f ...string) *Query {
    q.Sql["Group"] = fmt.Sprintf(" GROUP BY %s ", strings.Join(f, ", "))
    return q
}

// Having8jm
func (q *Query) Having(qs ...*Query) *Query {
    q.Sql["Having"] = fmt.Sprintf(" HAVING %s ", q.sqlCond)
    q.sqlCond = ""
    q.current = "Having"
    return q
}

// order by
func (q *Query) orderBy(sort string, f ...string) *Query {
    q.Sql["Order"] = fmt.Sprintf(" ORDER BY %s %s ", strings.Join(f, ", "), sort)
    q.current = "Order"
    return q
}

// Order ASC
func (q *Query) OrderAsc(f ...string) *Query {
    return q.orderBy("ASC", f...)
}

// Order DESC
func (q *Query) OrderDesc(f ...string) *Query {
    return q.orderBy("DESC", f...)
}

// Limit
func (q *Query) Limit(offset, rows int64) *Query {
    q.Sql["Limit"] = fmt.Sprintf(" LIMIT %d, %d ", offset, rows)
    q.current = "Limit"
    return q
}

// Delete
func (q *Query) DeleteFrom(tb string) *Query {
    q.Type = QueryDelete
    q.Sql["Delete"] = fmt.Sprintf(" DELETE FROM %s ", tb)
    q.current = "Delete"
    return q
}

// Update
func (q *Query) Update(tb string) *Query {
    q.Type = QueryUpdate
    q.Sql["Update"] = fmt.Sprintf(" UPDATE %s ", tb)
    q.current = "Update"
    return q
}

// Set(Update)
func (q *Query) Set(f string, v interface{}) *Query {
    str, ok := q.Sql["Set"]
    if ok && len(str) > 0 {
        q.Sql["Set"] += fmt.Sprintf(" , %s = ? ", f)
    } else {
        q.Sql["Set"] = fmt.Sprintf(" SET %s = ? ", f)
    }
    q.Args = append(q.Args, v)
    q.current = "Set"
    return q
}

// Insert
func (q *Query) InsertInto(tb string) *Query {
    q.Type = QueryInsert
    q.Sql["Insert"] = fmt.Sprintf(" INSERT INTO %s ", tb)
    q.current = "Insert"
    return q
}

// Fields(Insert)
func (q *Query) Fields(f ...string) *Query {
    q.Sql["Fields"] = fmt.Sprintf(" (%s) ", strings.Join(f, ", "))
    q.current = "Fields"
    return q
}

// Values(Insert)
func (q *Query) Values(vs ...interface{}) *Query {
    // Placeholder
    ph := make([]string, len(vs))
    for i, v := range vs {
        ph[i] = "?"
        q.Args = append(q.Args, v)
    }

    str, ok := q.Sql["Values"]
    if ok && len(str) > 0 {
        q.Sql["Values"] += fmt.Sprintf(" ,(%s) ", strings.Join(ph, ", "))
    } else {
        q.Sql["Values"] = fmt.Sprintf(" VALUES(%s) ", strings.Join(ph, ", "))
    }
    q.current = "Values"
    return q
}

// Connect all sql part to a corect sql string.
func (q *Query) ToString() string {
    str := ""
    nodes := queryNodes[q.Type]

    for _, node := range nodes {
        str += q.Sql[node]
    }

    log.Printf("%s\n", str)
    log.Printf("%#v\n", q.Args)
    return str
}



// Parse map data to insert SQL
func (q *Query) mapToInsert(d map[string]interface{}) {
    f := make([]string, 0)
    ph := make([]string, 0)
    for k, v := range d {
        f = append(f, k)
        ph = append(ph, "?")
        q.Args = append(q.Args, v)
    }
    q.Sql["Fields"] = fmt.Sprintf(" (%s) ", strings.Join(f, ", "))
    q.Sql["Values"] = fmt.Sprintf(" VALUES(%s) ", strings.Join(ph, ", "))
}

// Parse map data to update SQL
func (q *Query) mapToUpdate(d map[string]interface{}) {
    i := 0
    for k, v := range d {
        if i == 0 {
            q.Sql["Set"] = fmt.Sprintf(" SET %s = ? ", k)
        } else {
            q.Sql["Set"] += fmt.Sprintf(" , %s = ? ", k)
        }
        q.Args = append(q.Args, v)
        i++
    }
}

// Parse map data to where SQL
func (q *Query) mapToWhere(d map[string]interface{}) {
    i := 0
    for k, v := range d {
        if i == 0 {
            q.Sql["Where"] = fmt.Sprintf(" WHERE %s = ? ", k)
        } else {
            q.Sql["Where"] += fmt.Sprintf(" AND %s = ? ", k)
        }
        q.Args = append(q.Args, v)
        i++
    }
}

// Exec
func (q *Query) Exec(d ...map[string]interface{}) (sql.Result, error) {
    if q.Server == nil {
        return nil, errors.New("DB config not found")
    }
    switch q.Type {
        case QueryInsert:
        if len(d) == 1 {
            q.mapToInsert(d[0])
        }
        case QueryUpdate:
        if len(d) >= 1 {
            q.mapToUpdate(d[0])
        }
        if len(d) == 2 {
            q.mapToWhere(d[1])
        }
        case QueryDelete:
        if len(d) == 1 {
            q.mapToWhere(d[0])
        }
    }
    return q.Server.Exec(q.ToString(), q.Args...)
}

// Row
func (q *Query) Row(ptr interface{}, d ...map[string]interface{}) error {
    if q.Server == nil {
        return errors.New("DB config not found")
    }
    if len(d) == 1 {
        q.mapToWhere(d[0])
    }
    return q.Server.Row(ptr, q.ToString(), q.Args...)
}

// Rows
func (q *Query) Rows(ptr interface{}, d ...map[string]interface{}) error {
    if q.Server == nil {
        return errors.New("DB config not found")
    }
    if len(d) == 1 {
        q.mapToWhere(d[0])
    }
    return q.Server.Rows(ptr, q.ToString(), q.Args...)
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
