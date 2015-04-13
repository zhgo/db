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

// Database list
var Connections map[string]*DB = make(map[string]*DB)

// DB instance
var dbObjects map[string]*sql.DB = make(map[string]*sql.DB)

// Connect to database
func (e *DB) connect() error {
    db, s := dbObjects[e.Name]
    if s == false || db == nil { //err := e.Module.DB.Ping(); err != nil
        // s := "%s:%s@tcp(%s:%d)/%s?charset=utf8"
        // dsnStr := fmt.Sprintf(s, dsn.User, dsn.Password, dsn.Server, dsn.Port, dsn.Database)
        // log.Printf("%#v\n", c)

        var err error
        dbObjects[e.Name], err = sql.Open(e.Type, e.DSN)

        if err != nil {
            return err
        }
    }

    return nil
}

// get one row.
func (e *DB) Row(ptr interface{}, sql string, args []interface{}) error {
    rows, err := e.rows(sql, args)
    if err != nil {
        log.Printf("%s\n", err)
        return err
    }

    defer rows.Close()

    columns, err := rows.Columns()
    if err != nil {
        log.Printf("%s\n", err)
        return err
    }

    columnsLen := len(columns)

    dataKind, scanPtr, scanArgs, err := scanVariable(ptr, columnsLen, false)
    if err != nil {
        log.Printf("%s\n", err)
        return err
    }

    var scanVals []interface{}
    if dataKind == reflect.Map || dataKind == reflect.Slice {
        scanVals = reflect.ValueOf(scanPtr).Elem().Interface().([]interface{})
    }

    //返回值
    val := reflect.Indirect(reflect.ValueOf(ptr))

    for rows.Next() {
        if err := rows.Scan(scanArgs...); err != nil {
            log.Printf("%s\n", err)
            return err
        }

        switch dataKind {
            case reflect.Struct: //struct
            val.Set(reflect.Indirect(reflect.ValueOf(scanPtr)))
            case reflect.Map: //map
            //非指针, 必须放在for中每次定义
            record := make(map[string]interface{}, columnsLen)
            for i, col := range scanVals {
                record[columns[i]] = typeAssertion(col)
            }
            val.Set(reflect.ValueOf(record))
            case reflect.Slice: //slice
            record := make([]interface{}, columnsLen)
            for i, col := range scanVals {
                record[i] = typeAssertion(col)
            }
            val.Set(reflect.ValueOf(record))
        }

    }

    if err = rows.Err(); err != nil {
        log.Printf("%s\n", err)
        return err
    }

    return nil
}

// get all rows
func (e *DB) Rows(ptr interface{}, sql string, args []interface{}) error {
    rows, err := e.rows(sql, args)
    if err != nil {
        log.Printf("%s\n", err)
        return err
    }

    defer rows.Close()

    columns, err := rows.Columns()
    if err != nil {
        log.Printf("%s\n", err)
        return err
    }

    columnsLen := len(columns)

    dataKind, scanPtr, scanArgs, err := scanVariable(ptr, columnsLen, true)
    if err != nil {
        log.Printf("%s\n", err)
        return err
    }

    var scanVals []interface{}
    if dataKind == reflect.Map || dataKind == reflect.Slice {
        scanVals = reflect.ValueOf(scanPtr).Elem().Interface().([]interface{})
    }

    //return data
    val := reflect.Indirect(reflect.ValueOf(ptr))

    for rows.Next() {
        if err := rows.Scan(scanArgs...); err != nil {
            log.Printf("%s\n", err)
            return err
        }

        switch dataKind {
            case reflect.Struct: //struct
            val.Set(reflect.Append(val, reflect.Indirect(reflect.ValueOf(scanPtr))))
            case reflect.Map: //map
            //非指针, 必须放在for中每次定义
            record := make(map[string]interface{}, columnsLen)
            for i, col := range scanVals {
                record[columns[i]] = typeAssertion(col)
            }
            val.Set(reflect.Append(val, reflect.ValueOf(record)))
            case reflect.Slice: //slice
            record := make([]interface{}, columnsLen)
            for i, col := range scanVals {
                record[i] = typeAssertion(col)
            }
            val.Set(reflect.Append(val, reflect.ValueOf(record)))
        }
    }

    if err = rows.Err(); err != nil {
        log.Printf("%s\n", err)
        return err
    }

    return nil
}

// Execute query, only return sql.Result
func (e *DB) Exec(sql string, args []interface{}) (sql.Result, error) {
    stmt, err := e.Prepare(sql)
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

// Execute query, return sql.Rows
func (e *DB) rows(sql string, args []interface{}) (*sql.Rows, error) {
    stmt, err := e.Prepare(sql)
    if err != nil {
        return nil, err
    }

    defer stmt.Close()

    rows, err := stmt.Query(args...)
    if err != nil {
        return nil, err
    }

    return rows, nil
}

// Prepare
func (e *DB) Prepare(sql string) (*sql.Stmt, error) {
    if err := e.connect(); err != nil {
        return nil, err
    }

    stmt, err := dbObjects[e.Name].Prepare(sql)
    if err != nil {
        return nil, err
    }

    return stmt, nil
}

