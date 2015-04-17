// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
    "testing"
    "io/ioutil"
)

func TestQueryMysql(t *testing.T) {
    server := Server{
        Name: "test-mysql",
        Type: "mysql",
        DSN: "root:@tcp(127.0.0.1:3306)/zhgo?charset=utf8"}


    // Load tb-mysql.sql
    b, err := ioutil.ReadFile("tb-mysql.sql")
    if err != nil {
        t.Fatalf("Read files failed (tb-mysql.sql): %v.\n", err)
    }


    // Drop table1
    query := NewQuery(&server)
    err = query.DropTable("table1")
    if err != nil {
        t.Fatalf("Drop table1 failed: %v.\n", err)
    }


    // Create table1
    query = NewQuery(&server)
    _, err = query.Server.Exec(string(b))
    if err != nil {
        t.Fatalf("Create table1 failed: %v.\n", err)
    }


    // Insert
    query = NewQuery(&server)
    query.InsertInto("table1")
    query.Fields("BirthYear", "Gender", "Nickname")
    query.Values("1980", "Male", "张三丰")
    r, err := query.Exec()
    if err != nil {
        t.Fatalf("Insert data to table1 failed: %v.\n", err)
    }
    lastInsertId, err := r.LastInsertId()
    if err != nil {
        t.Fatalf("Insert data to table1 failed: %v.\n", err)
    }
    if lastInsertId != 1000000 {
        t.Fatalf("Insert data to table1 failed: LastInsertId error.\n")
    }


    // Insert confirm
    d := make(map[string]interface{})
    query = NewQuery(&server)
    err = query.Select("*").From("table1").Where(query.Eq("UserID", "1000000")).Row(&d)
    if err != nil {
        t.Fatalf("Select table1 failed: %v.\n", err)
    }
    if d["BirthYear"] != int64(1980) {
        t.Fatalf("table1 data error (BirthYear): %v.\n", d["BirthYear"])
    }
    if d["Gender"] != "Male" {
        t.Fatalf("table1 data error (Gender): %v.\n", d["Gender"])
    }
    if d["Nickname"] != "张三丰" {
        t.Fatalf("table1 data error (Nickname): %v.\n", d["Nickname"])
    }


    // Update
    query = NewQuery(&server)
    r, err = query.Update("table1").Set("BirthYear", "1982").Set("Gender", "Female").Set("Nickname", "Bob").Where(query.Eq("UserID", "1000000")).Exec()
    if err != nil {
        t.Fatalf("Update table1 failed: %v.\n", err)
    }
    rowsAffected, err := r.RowsAffected()
    if err != nil {
        t.Fatalf("Update table1 failed: %v.\n", err)
    }
    if rowsAffected <= 0 {
        t.Fatalf("Update table1 failed.\n")
    }

    // Update confirm
    d = make(map[string]interface{})
    query = NewQuery(&server)
    err = query.Select("*").From("table1").Where(query.Eq("UserID", "1000000")).Row(&d)
    if err != nil {
        t.Fatalf("Select table1 failed: %v.\n", err)
    }
    if d["BirthYear"] != int64(1982) {
        t.Fatalf("table1 data error (BirthYear): %v.\n", d["BirthYear"])
    }
    if d["Gender"] != "Female" {
        t.Fatalf("table1 data error (Gender): %v.\n", d["Gender"])
    }
    if d["Nickname"] != "Bob" {
        t.Fatalf("table1 data error (Nickname): %v.\n", d["Nickname"])
    }
}

func TestQuerySqlite3(t *testing.T) {
    server := NewServer("test-sqlite3", "sqlite3", "sqlite3.db")


    // Load tb-sqlite3.sql
    b, err := ioutil.ReadFile("tb-sqlite3.sql")
    if err != nil {
        t.Fatalf("Read files failed (tb-sqlite3.sql): %v.\n", err)
    }


    // Drop table1
    query := NewQuery(server)
    err = query.DropTable("table1")
    if err != nil {
        t.Fatalf("Drop table1 failed: %v.\n", err)
    }


    // Create table1
    query = NewQuery(server)
    _, err = query.Server.Exec(string(b))
    if err != nil {
        t.Fatalf("Create table1 failed: %v.\n", err)
    }


    // Insert
    query = NewQuery(server)
    query.InsertInto("table1")
    query.Fields("UserID", "CreationTime", "BirthYear", "Gender", "Nickname")
    query.Values("1000000", "1429091207", "1980", "Male", "张三丰")
    r, err := query.Exec()
    if err != nil {
        t.Fatalf("Insert data to table1 failed: %v.\n", err)
    }
    lastInsertId, err := r.LastInsertId()
    if err != nil {
        t.Fatalf("Insert data to table1 failed: %v.\n", err)
    }
    if lastInsertId != 1000000 {
        t.Fatalf("Insert data to table1 failed: LastInsertId error.\n")
    }


    // Insert confirm
    d := make(map[string]interface{})
    query = NewQuery(server)
    err = query.Select("*").From("table1").Where(query.Eq("UserID", "1000000")).Row(&d)
    if err != nil {
        t.Fatalf("Select table1 failed: %v.\n", err)
    }
    if d["CreationTime"] != int64(1429091207) {
        t.Fatalf("table1 data error (CreationTime): %v.\n", d["CreationTime"])
    }
    if d["BirthYear"] != int64(1980) {
        t.Fatalf("table1 data error (BirthYear): %v.\n", d["BirthYear"])
    }
    if d["Gender"] != "Male" {
        t.Fatalf("table1 data error (Gender): %v.\n", d["Gender"])
    }
    if d["Nickname"] != "张三丰" {
        t.Fatalf("table1 data error (Nickname): %v.\n", d["Nickname"])
    }


    // Update
    query = NewQuery(server)
    r, err = query.Update("table1").Set("BirthYear", "1982").Set("Gender", "Female").Set("Nickname", "Bob").Where(query.Eq("UserID", "1000000")).Exec()
    if err != nil {
        t.Fatalf("Update table1 failed: %v.\n", err)
    }
    rowsAffected, err := r.RowsAffected()
    if err != nil {
        t.Fatalf("Update table1 failed: %v.\n", err)
    }
    if rowsAffected <= 0 {
        t.Fatalf("Update table1 failed.\n")
    }

    // Update confirm
    d = make(map[string]interface{})
    query = NewQuery(server)
    err = query.Select("*").From("table1").Where(query.Eq("UserID", "1000000")).Row(&d)
    if err != nil {
        t.Fatalf("Select table1 failed: %v.\n", err)
    }
    if d["BirthYear"] != int64(1982) {
        t.Fatalf("table1 data error (BirthYear): %v.\n", d["BirthYear"])
    }
    if d["Gender"] != "Female" {
        t.Fatalf("table1 data error (Gender): %v.\n", d["Gender"])
    }
    if d["Nickname"] != "Bob" {
        t.Fatalf("table1 data error (Nickname): %v.\n", d["Nickname"])
    }
}