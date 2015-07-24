// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	_ "github.com/zhgo/mysql"
	_ "github.com/zhgo/postgresql"
	_ "github.com/zhgo/sqlite/sqlite3"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Server struct
type Server struct {
	// Name
	Name string

	// Database type: mysql postgresql or sqlite3
	Type string

	// Data Source Name
	DSN string

	// Follow
	Follow string
}

// Server instance
var dbObjects map[string]*sql.DB = make(map[string]*sql.DB)

// Execute query, only return sql.Result
func (e *Server) Exec(sql string, args ...interface{}) (sql.Result, error) {
	sql, args = e.parseSQL(sql, args)
	stmt, err := e.prepare(sql)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Get row.
func (e *Server) Row(ptr interface{}, sql string, args ...interface{}) error {
	rows, columns, err := e.rows(sql, args)
	if err != nil {
		log.Printf("%s\n", err)
		return err
	}

	defer rows.Close()

	columnsLen := len(columns)

	kind, ptrRow, scan, err := scanVariables(ptr, columnsLen, false)
	if err != nil {
		log.Printf("%s\n", err)
		return err
	}

	// Return data
	val := reflect.ValueOf(ptr).Elem()

	if rows.Next() {
		if err := rows.Scan(scan...); err != nil {
			log.Printf("%s\n", err)
			return err
		}

		switch kind {
		case reflect.Struct: // struct
			val.Set(reflect.ValueOf(ptrRow).Elem())

		case reflect.Map: //map
			row := make(map[string]interface{}, columnsLen)
			for i := 0; i < columnsLen; i++ {
				row[columns[i]] = typeAssertion(*(scan[i].(*interface{})))
			}
			val.Set(reflect.ValueOf(row))

		case reflect.Slice: //slice
			row := make([]interface{}, columnsLen)
			for i := 0; i < columnsLen; i++ {
				row[i] = typeAssertion(*(scan[i].(*interface{})))
			}
			val.Set(reflect.ValueOf(row))
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("%s\n", err)
		return err
	}

	return nil
}

// Get all rows
func (e *Server) Rows(ptr interface{}, sql string, args ...interface{}) error {
	rows, columns, err := e.rows(sql, args)
	if err != nil {
		log.Printf("%s\n", err)
		return err
	}

	defer rows.Close()

	columnsLen := len(columns)

	kind, ptrRow, scan, err := scanVariables(ptr, columnsLen, true)
	if err != nil {
		log.Printf("%s\n", err)
		return err
	}

	//return data
	val := reflect.ValueOf(ptr).Elem()

	for rows.Next() {
		if err := rows.Scan(scan...); err != nil {
			log.Printf("%s\n", err)
			return err
		}

		switch kind {
		case reflect.Struct: // struct
			val.Set(reflect.Append(val, reflect.ValueOf(ptrRow).Elem()))

		case reflect.Map: // map
			row := make(map[string]interface{}, columnsLen)
			for i := 0; i < columnsLen; i++ {
				row[columns[i]] = typeAssertion(*(scan[i].(*interface{})))
			}
			val.Set(reflect.Append(val, reflect.ValueOf(row)))

		case reflect.Slice: // slice
			row := make([]interface{}, columnsLen)
			for i := 0; i < columnsLen; i++ {
				row[i] = typeAssertion(*(scan[i].(*interface{})))
			}
			val.Set(reflect.Append(val, reflect.ValueOf(row)))
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("%s\n", err)
		return err
	}

	return nil
}

// New query
func (e *Server) NewQuery() *Query {
	return NewQuery(e)
}

// Insert into
func (e *Server) InsertInto(tb string) *Query {
	return NewQuery(e).InsertInto(tb)
}

// Update
func (e *Server) Update(tb string) *Query {
	return NewQuery(e).Update(tb)
}

// Delete from
func (e *Server) DeleteFrom(tb string) *Query {
	return NewQuery(e).DeleteFrom(tb)
}

// Select
func (e *Server) Select(f ...string) *Query {
	return NewQuery(e).Select(f...)
}

// Execute query, return sql.Rows, rows.Columns
func (e *Server) rows(sql string, args []interface{}) (*sql.Rows, []string, error) {
	sql, args = e.parseSQL(sql, args)
	stmt, err := e.prepare(sql)
	if err != nil {
		return nil, nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	return rows, columns, nil
}

// Prepare SQL
func (e *Server) prepare(sql string) (*sql.Stmt, error) {
	if err := e.connect(); err != nil {
		return nil, err
	}

	stmt, err := dbObjects[e.Name].Prepare(sql)
	if err != nil {
		return nil, err
	}

	return stmt, nil
}

// Connect to database
func (e *Server) connect() error {
	db, s := dbObjects[e.Name]
	// err := db.Ping(); err != nil
	if s == false || db == nil {
		var err error
		dbObjects[e.Name], err = sql.Open(e.Type, e.DSN)
		if err != nil {
			return err
		}
	}

	return nil
}

// sql compatibility
func (e *Server) parseSQL(str string, args []interface{}) (string, []interface{}) {
	switch e.Type {
	case "postgres":
		return str, args
	case "mysql":
		str = e.parseQuotes(str)
		return e.parseParameters(str, args)
	case "sqlite3":
		return e.parseParameters(str, args)
	}
	return str, args
}

// " to `
func (e *Server) parseQuotes(str string) string {
	return strings.Replace(str, `"`, "`", -1)
}

// $1, $2, $3 to ?, ?, ?
func (e *Server) parseParameters(str string, args []interface{}) (string, []interface{}) {
	re := regexp.MustCompile(`\$(\d+)`)
	sli := re.FindAllStringSubmatch(str, -1)
	newArgs := make([]interface{}, len(sli))
	for i, v := range sli {
		vi, err := strconv.ParseInt(v[1], 10, 0)
		if err != nil {
			log.Printf("%s\n", err)
		}
		newArgs[i] = args[vi-1]
	}
	return re.ReplaceAllString(str, "?"), newArgs
}

// New Server
func NewServer(name string, typ string, dsn string) *Server {
	return &Server{Name: typ + name, Type: typ, DSN: dsn}
}

// New Server, alias of NewServer.
func Connect(name string, typ string, dsn string) *Server {
	return NewServer(name, typ, dsn)
}
