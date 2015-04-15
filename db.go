// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/mxk/go-sqlite/sqlite3"
    "log"
    "reflect"
)

// Database struct
type DB struct {
    // Name
    Name string

    // Database type: mysql or sqlite3
    Type string

    // Data Source Name
    DSN string

    // profiling
    Profiling bool

    // Follow
    Follow string
}

// DB instance
var dbObjects map[string]*sql.DB = make(map[string]*sql.DB)

// Execute query, only return sql.Result
func (e *DB) Exec(sql string, args []interface{}) (sql.Result, error) {
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
func (e *DB) Row(ptr interface{}, sql string, args []interface{}) error {
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
func (e *DB) Rows(ptr interface{}, sql string, args []interface{}) error {
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

// Execute query, return sql.Rows, rows.Columns
func (e *DB) rows(sql string, args []interface{}) (*sql.Rows, []string, error) {
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
func (e *DB) prepare(sql string) (*sql.Stmt, error) {
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
func (e *DB) connect() error {
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

