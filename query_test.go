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

    query := NewQuery(&server)


    // Load table1.sql
    b, err := ioutil.ReadFile("tb-mysql.sql")
    if err != nil {
        t.Fatalf("Read files failed (tb-mysql.sql): %v.\n", err)
    }


    // Drop table1
    err = query.DropTable("table1")
    if err != nil {
        t.Fatalf("Drop table1 failed: %v.\n", err)
    }


    // Create table1
    _, err = server.Exec(string(b))
    if err != nil {
        t.Fatalf("Create table1 failed: %v.\n", err)
    }


    // Insert
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
}