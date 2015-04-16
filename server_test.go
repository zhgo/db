// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
    "database/sql"
    "io/ioutil"
    "testing"
)

/*
CREATE USER 'zhgo'@'localhost' IDENTIFIED BY 'zhgo';
CREATE DATABASE `zhgo` CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, ALTER, LOCK TABLES ON `zhgo`.* TO 'zhgo'@'localhost';
*/

func TestDBMysql(t *testing.T) {
    var q string
    var r sql.Result

    db := Server{
        Name: "test-mysql",
        Type: "mysql",
        DSN: "root:@tcp(127.0.0.1:3306)/zhgo?charset=utf8",
        Profiling: false,
        Follow: ""}


    // Load table1.sql
    b, err := ioutil.ReadFile("tb-mysql.sql")
    if err != nil {
        t.Fatalf("Read files failed (tb-mysql.sql): %v.\n", err)
    }


    // Drop table1
    q = "DROP TABLE IF EXISTS table1"
    _, err = db.Exec(q, nil)
    if err != nil {
        t.Fatalf("Drop table1 failed: %v.\n", err)
    }


    // Create table1
    _, err = db.Exec(string(b), nil)
    if err != nil {
        t.Fatalf("Create table1 failed: %v.\n", err)
    }


    // Insert
    q = "INSERT INTO table1(BirthYear, Gender, Nickname) VALUES('1980', 'Male', '张三丰')"
    r, err = db.Exec(q, nil)
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
    q = "SELECT * FROM table1 WHERE UserID = ?"
    p := []interface{}{1000000}
    err = db.Row(&d, q, p)
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
    q = "UPDATE table1 SET BirthYear = ?, Gender = ?, Nickname = ? WHERE UserID = '1000000'"
    p = []interface{}{1982, "Female", "Bob"}
    r, err = db.Exec(q, p)
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
    q = "SELECT * FROM table1 WHERE UserID = ?"
    p = []interface{}{1000000}
    err = db.Row(&d, q, p)
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

func TestDBSqlite3(t *testing.T) {
    var q string
    var r sql.Result

    db := Server{
        Name: "test-sqlite3",
        Type: "sqlite3",
        DSN: "sqlite3.db",
        Profiling: false,
        Follow: ""}


    // Load table1.sql
    b, err := ioutil.ReadFile("tb-sqlite3.sql")
    if err != nil {
        t.Fatalf("Read files failed (tb-sqlite3.sql): %v.\n", err)
    }


    // Drop table1
    q = "DROP TABLE IF EXISTS table1"
    _, err = db.Exec(q, nil)
    if err != nil {
        t.Fatalf("Drop table1 failed: %v.\n", err)
    }


    // Create table1
    _, err = db.Exec(string(b), nil)
    if err != nil {
        t.Fatalf("Create table1 failed: %v.\n", err)
    }


    // Insert
    q = "INSERT INTO table1(UserID, CreationTime, BirthYear, Gender, Nickname) VALUES(1000000, '1429091207', '1980', 'Male', '张三丰')"
    r, err = db.Exec(q, nil)
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
    q = "SELECT * FROM table1 WHERE UserID = ?"
    p := []interface{}{1000000}
    err = db.Row(&d, q, p)
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
    q = "UPDATE table1 SET BirthYear = ?, Gender = ?, Nickname = ? WHERE UserID = '1000000'"
    p = []interface{}{1982, "Female", "Bob"}
    r, err = db.Exec(q, p)
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
    q = "SELECT * FROM table1 WHERE UserID = ?"
    p = []interface{}{1000000}
    err = db.Row(&d, q, p)
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
