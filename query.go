// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"errors"
	"fmt"
	"log"
	"strings"
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
	QueryInsert: []string{"Insert", "Fields", "Values", "Returning"},
	QueryUpdate: []string{"Update", "Set", "Where", "Order", "Limit"},
	QueryDelete: []string{"Delete", "Where"},
	QuerySelect: []string{"Select", "From", "Join", "Where", "Group", "Having", "Order", "Limit", "ForUpdate"}}

// Query struct
type Query struct {
	// Server
	Server *Server

	// Query type: Insert, Update, Delete, Select
	Type uint

	// Table primary field
	Primary string

	// Sql
	Sql map[string]string

	// Condition SQL
	SqlCond []string

	//args for Sql
	Args []interface{}

	//args index
	ArgIndex int

	// Current Sql node
	current string
}

// A Result summarizes an executed SQL command.
type Result struct {
	// LastInsertId returns the integer generated by the database
	// in response to a command. Typically this will be from an
	// "auto increment" column when inserting a new row. Not all
	// databases support this feature, and the syntax of such
	// statements varies.
	LastInsertId int64

	// RowsAffected returns the number of rows affected by an
	// update, insert, or delete. Not every database or database
	// driver may support this.
	RowsAffected int64
}

// Condition struct
type Condition struct {
	Eq   map[string]string   `json:"eq"`
	Ge   map[string]string   `json:"ge"`
	Gt   map[string]string   `json:"gt"`
	Le   map[string]string   `json:"le"`
	Lt   map[string]string   `json:"lt"`
	Ne   map[string]string   `json:"ne"`
	Like map[string]string   `json:"like"`
	In   map[string][]string `json:"in"`
}

// Equal
func (q *Query) Eq(f string, v interface{}) string {
	return q.condition("", f, "=", v)
}

// Greater than or equal
func (q *Query) Ge(f string, v interface{}) string {
	return q.condition("", f, ">=", v)
}

// Greater than
func (q *Query) Gt(f string, v interface{}) string {
	return q.condition("", f, ">", v)
}

// Less than or equal
func (q *Query) Le(f string, v interface{}) string {
	return q.condition("", f, "<=", v)
}

// Less than
func (q *Query) Lt(f string, v interface{}) string {
	return q.condition("", f, "<", v)
}

// Not equal
func (q *Query) Ne(f string, v interface{}) string {
	return q.condition("", f, "<>", v)
}

// Like: Simple pattern matching
func (q *Query) Like(f string, v interface{}) string {
	return q.condition("", f, "LIKE", v)
}

// In: Check whether a value is within a set of values
func (q *Query) In(f string, v ...interface{}) string {
	return q.conditionIn("", f, v...)
}

// And
func (q *Query) And(qs ...string) string {
	return fmt.Sprintf(" AND (%s) ", strings.Join(qs, " "))
}

// and
func (q *Query) and(f string, co string, v interface{}) string {
	return q.condition("AND", f, co, v)
}

// And (Equal)
func (q *Query) AndEq(f string, v interface{}) string {
	return q.and(f, "=", v)
}

// And (Greater than or equal)
func (q *Query) AndGe(f string, v interface{}) string {
	return q.and(f, ">=", v)
}

// And (Greater than)
func (q *Query) AndGt(f string, v interface{}) string {
	return q.and(f, ">", v)
}

// And (Less than or equal)
func (q *Query) AndLe(f string, v interface{}) string {
	return q.and(f, "<=", v)
}

// And (Less than)
func (q *Query) AndLt(f string, v interface{}) string {
	return q.and(f, "<", v)
}

// And (Not equal)
func (q *Query) AndNe(f string, v interface{}) string {
	return q.and(f, "<>", v)
}

// And (Like: Simple pattern matching)
func (q *Query) AndLike(f string, v interface{}) string {
	return q.and(f, "LIKE", v)
}

// And (In: Check whether a value is within a set of values)
func (q *Query) AndIn(f string, v ...interface{}) string {
	return q.conditionIn("AND", f, v...)
}

// Or
func (q *Query) Or(qs ...string) string {
	return fmt.Sprintf(" OR (%s) ", strings.Join(qs, " "))
}

// or
func (q *Query) or(f string, co string, v interface{}) string {
	return q.condition("OR", f, co, v)
}

// Or (Equal)
func (q *Query) OrEq(f string, v interface{}) string {
	return q.or(f, "=", v)
}

// Or (Greater than or equal)
func (q *Query) OrGe(f string, v interface{}) string {
	return q.or(f, ">=", v)
}

// Or (Greater than)
func (q *Query) OrGt(f string, v interface{}) string {
	return q.or(f, ">", v)
}

// Or (Less than or equal)
func (q *Query) OrLe(f string, v interface{}) string {
	return q.or(f, "<=", v)
}

// Or (Less than)
func (q *Query) OrLt(f string, v interface{}) string {
	return q.or(f, "<", v)
}

// Or (Not equal)
func (q *Query) OrNe(f string, v interface{}) string {
	return q.or(f, "<>", v)
}

// Or (Like: Simple pattern matching)
func (q *Query) OrLike(f string, v interface{}) string {
	return q.or(f, "LIKE", v)
}

// Or (In: Check whether a value is within a set of values)
func (q *Query) OrIn(f string, v ...interface{}) string {
	return q.conditionIn("OR", f, v...)
}

// l Logical
// f Field
// co Comparison Operators
// v Value
func (q *Query) condition(l string, f string, co string, v interface{}) string {
	return fmt.Sprintf(" %s %s %s %s ", l, q.quoteField(f), co, q.placeholder(v))
}

// l Logical
// f Field
// vs Values
func (q *Query) conditionIn(l string, f string, vs ...interface{}) string {
	ph := make([]string, len(vs)) // Placeholder
	for i, v := range vs {
		ph[i] = q.placeholder(v)
	}
	return fmt.Sprintf(" %s %s IN (%s) ", l, q.quoteField(f), strings.Join(ph, ", "))
}

// Join fields
func (q *Query) joinFields(fs []string) string {
	return strings.Join(q.quoteFields(fs), ", ")
}

// Quote fileds
func (q *Query) quoteFields(f []string) []string {
	// Slices are passed by reference, copy f to nf.
	nf := make([]string, len(f))
	for i, v := range f {
		nf[i] = q.quoteField(v)
	}
	return nf
}

// Quote filed
func (q *Query) quoteField(f string) string {
	fs := strings.Trim(f, " \r\n\t")
	if fs == "*" || fs == "1" {
		return f
	}
	return fmt.Sprintf("\"%s\"", strings.Replace(f, ".", "\".\"", -1))
}

// Placeholder
func (q *Query) placeholder(v interface{}) string {
	q.ArgIndex++
	q.Args = append(q.Args, v)
	return fmt.Sprintf("$%d", q.ArgIndex)
}

// Insert
func (q *Query) InsertInto(tb string) *Query {
	q.Type = QueryInsert
	q.Sql["Insert"] = fmt.Sprintf(" INSERT INTO %s ", q.quoteField(tb))
	q.current = "Insert"
	return q
}

// Set primary field
func (q *Query) SetPrimary(p string) {
	q.Primary = p
}

// Fields(Insert)
func (q *Query) Fields(f ...string) *Query {
	q.Sql["Fields"] = fmt.Sprintf(" (%s) ", q.joinFields(f))
	q.current = "Fields"
	return q
}

// Values(Insert)
func (q *Query) Values(vs ...interface{}) *Query {
	// Placeholder
	ph := make([]string, len(vs))
	for i, v := range vs {
		ph[i] = q.placeholder(v)
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

// Update
func (q *Query) Update(tb string) *Query {
	q.Type = QueryUpdate
	q.Sql["Update"] = fmt.Sprintf(" UPDATE %s ", q.quoteField(tb))
	q.current = "Update"
	return q
}

// Set(Update)
func (q *Query) Set(f string, v interface{}) *Query {
	str, ok := q.Sql["Set"]
	if ok && len(str) > 0 {
		q.Sql["Set"] += fmt.Sprintf(" , %s = %s ", q.quoteField(f), q.placeholder(v))
	} else {
		q.Sql["Set"] = fmt.Sprintf(" SET %s = %s ", q.quoteField(f), q.placeholder(v))
	}
	q.current = "Set"
	return q
}

// Delete
func (q *Query) DeleteFrom(tb string) *Query {
	q.Type = QueryDelete
	q.Sql["Delete"] = fmt.Sprintf(" DELETE FROM %s ", q.quoteField(tb))
	q.current = "Delete"
	return q
}

// Select fields
func (q *Query) Select(f ...string) *Query {
	q.Type = QuerySelect
	if len(f) == 0 {
		// Warring!
		// If use struct type as return data type, must keep have same sequence
		// about struct fileds and table fileds.
		q.Sql["Select"] = " SELECT *"
	} else {
		q.Sql["Select"] = fmt.Sprintf(" SELECT %s ", q.joinFields(f))
	}
	q.current = "Select"
	return q
}

// From
func (q *Query) From(tb string) *Query {
	q.Sql["From"] = fmt.Sprintf(" FROM %s ", q.quoteField(tb))
	q.current = "From"
	return q
}

// join
func (q *Query) join(p string, tb string) *Query {
	q.Sql["Join"] += fmt.Sprintf(" %s JOIN %s ", p, q.quoteField(tb))
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
func (q *Query) On(qs ...string) *Query {
	q.Sql["Join"] += fmt.Sprintf(" ON %s ", strings.Join(qs, " "))
	q.SqlCond = make([]string, 0)
	q.current = "Join"
	return q
}

// Where
func (q *Query) Where(qs ...string) *Query {
	q.Sql["Where"] = fmt.Sprintf(" WHERE %s ", strings.Join(qs, " "))
	q.SqlCond = make([]string, 0)
	q.current = "Where"
	return q
}

// Group By
func (q *Query) GroupBy(f ...string) *Query {
	q.Sql["Group"] = fmt.Sprintf(" GROUP BY %s ", q.joinFields(f))
	return q
}

// Having
func (q *Query) Having(qs ...string) *Query {
	q.Sql["Having"] = fmt.Sprintf(" HAVING %s ", strings.Join(qs, " "))
	q.SqlCond = make([]string, 0)
	q.current = "Having"
	return q
}

// order by
func (q *Query) orderBy(sort string, f ...string) *Query {
	q.Sql["Order"] = fmt.Sprintf(" ORDER BY %s %s ", q.joinFields(f), sort)
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

// Parse
func (q *Query) Parse(c Condition) *Query {
	conds := []string{q.Eq("1", "1")}

	for k, v := range c.Eq {
		conds = append(conds, q.AndEq(k, v))
	}

	for k, v := range c.Ge {
		conds = append(conds, q.AndGe(k, v))
	}

	for k, v := range c.Gt {
		conds = append(conds, q.AndGt(k, v))
	}

	for k, v := range c.Le {
		conds = append(conds, q.AndLe(k, v))
	}

	for k, v := range c.Lt {
		conds = append(conds, q.AndLt(k, v))
	}

	for k, v := range c.Ne {
		conds = append(conds, q.AndNe(k, v))
	}

	for k, v := range c.Like {
		conds = append(conds, q.AndLike(k, v))
	}

	for k, v := range c.In {
		conds = append(conds, q.AndIn(k, v))
	}

	return q.Where(conds...)
}

// Connect all sql part to a corect sql string.
func (q *Query) ToString() string {
	str := ""
	nodes := queryNodes[q.Type]

	for _, node := range nodes {
		str += q.Sql[node]
	}

	if Env < 2 {
		log.Printf("%s\n", str)
		log.Printf("%#v\n", q.Args)
	} else {
		// FIXME: Write to file or redis or solr.
	}

	return str
}

// Parse map data to insert SQL
func (q *Query) mapToInsert(d Item) {
	f := make([]string, 0)
	ph := make([]string, 0)
	for k, v := range d {
		f = append(f, k)
		ph = append(ph, q.placeholder(v))
	}
	q.Sql["Fields"] = fmt.Sprintf(" (%s) ", q.joinFields(f))
	q.Sql["Values"] = fmt.Sprintf(" VALUES(%s) ", strings.Join(ph, ", "))
}

// Parse map data to update SQL
func (q *Query) mapToUpdate(d Item) {
	i := 0
	for k, v := range d {
		if i == 0 {
			q.Sql["Set"] = fmt.Sprintf(" SET %s = %s ", q.quoteField(k), q.placeholder(v))
		} else {
			q.Sql["Set"] += fmt.Sprintf(" , %s = %s ", q.quoteField(k), q.placeholder(v))
		}
		i++
	}
}

// Parse map data to where SQL
func (q *Query) mapToWhere(d Where) {
	i := 0
	for k, v := range d {
		if i == 0 {
			q.Sql["Where"] = fmt.Sprintf(" WHERE %s = %s ", q.quoteField(k), q.placeholder(v))
		} else {
			q.Sql["Where"] += fmt.Sprintf(" AND %s = %s ", q.quoteField(k), q.placeholder(v))
		}
		i++
	}
}

// Exec
func (q *Query) Exec(d ...map[string]interface{}) (Result, error) {
	re := Result{}
	if q.Server == nil {
		return re, errors.New("DB config not found")
	}

	switch q.Type {
	case QueryInsert:
		if len(d) == 1 {
			q.mapToInsert(d[0])
		}

		// https://github.com/lib/pq/issues/24
		if q.Server.Type == "postgres" {
			q.Sql["Returning"] = fmt.Sprintf("RETURNING %s", q.quoteField(q.Primary))
			row := make(Item)
			err := q.Server.Row(&row, q.ToString(), q.Args...)
			if err != nil {
				return re, err
			}
			lastInsertId, ok := row[q.Primary]
			if !ok {
				return re, errors.New(fmt.Sprintf("no LastInsertId available: %#v", row))
			}
			re.LastInsertId = lastInsertId.(int64)
			return re, nil
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

	r, err := q.Server.Exec(q.ToString(), q.Args...)
	if err != nil {
		return re, err
	}
	lastInsertId, err := r.LastInsertId()
	if err != nil && q.Server.Type != "postgres" {
		return re, err
	}
	re.LastInsertId = lastInsertId
	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return re, err
	}
	re.RowsAffected = rowsAffected
	return re, nil
}

// Row
func (q *Query) Row(ptr interface{}, d ...Where) error {
	if q.Server == nil {
		return errors.New("DB config not found")
	}
	if len(d) == 1 {
		q.mapToWhere(d[0])
	}
	return q.Server.Row(ptr, q.ToString(), q.Args...)
}

// Rows
func (q *Query) Rows(ptr interface{}, d ...Where) error {
	if q.Server == nil {
		return errors.New("DB config not found")
	}
	if len(d) == 1 {
		q.mapToWhere(d[0])
	}
	return q.Server.Rows(ptr, q.ToString(), q.Args...)
}

// New Query object
func NewQuery(server *Server) *Query {
	query := Query{Server: server}
	query.Sql = make(map[string]string)
	query.Args = make([]interface{}, 0)
	return &query
}
