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
    err = query.Select("*").From("table1").Where("UserID", "=", "1000000").Row(&d)
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

}