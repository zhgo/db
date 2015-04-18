// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
    "io/ioutil"
    "testing"
    "fmt"
    "strings"
)

/*
CREATE USER 'zhgo'@'localhost' IDENTIFIED BY 'zhgo';
CREATE DATABASE `zhgo` CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, ALTER, LOCK TABLES ON `zhgo`.* TO 'zhgo'@'localhost';
*/

type ServerTest struct {
    // Server
    Server *Server
}

func (st *ServerTest) Init(t *testing.T) {
    // Load sql file
    p := fmt.Sprintf("tb-%s.sql", st.Server.Type)
    b, err := ioutil.ReadFile(p)
    if err != nil {
        t.Fatalf("[%s]: %v\n", st.Server.Type, err)
    }
    initSQL := strings.Split(string(b), ";")

    // Run Init SQL
    for _, v := range initSQL {
        if strings.Trim(v, "\r\n \t") != "" {
            _, err := st.Server.Exec(v)
            if err != nil {
                t.Fatalf("[%s]: %v\n", st.Server.Type, err)
            }
        }
    }
}

func (st *ServerTest) Load(t *testing.T) {

}

func (st *ServerTest) Start(t *testing.T) {
    st.Insert(t)
    st.Update(t)
    st.Rows(t)
    st.Delete(t)
}

func (st *ServerTest) Insert(t *testing.T) {
    // Insert
    q := "INSERT INTO passport_user(UserID, CreationTime, BirthYear, Gender, Nickname) VALUES('1000000', '2015-01-17 00:00:00', '1980', 'Male', '肯·汤普逊')"
    r, err := st.Server.Exec(q)
    if err != nil {
        t.Fatalf("[%s]: %v\n", st.Server.Type, err)
    }
    lastInsertId, err := r.LastInsertId()
    if err != nil {
        t.Fatalf("[%s]: %v\n", st.Server.Type, err)
    }
    if lastInsertId != 1000000 {
        t.Fatalf("[%s]: %v\n", st.Server.Type, lastInsertId)
    }


    // Insert confirm
    d := make(map[string]interface{})
    q = "SELECT * FROM passport_user WHERE UserID = ?"
    err = st.Server.Row(&d, q, 1000000)
    if err != nil {
        t.Fatalf("[%s]: %v\n", st.Server.Type, err)
    }
    st.dataValidation(t, d["CreationTime"], "2015-01-17 00:00:00")
    st.dataValidation(t, d["BirthYear"], int64(1980))
    st.dataValidation(t, d["Gender"], "Male")
    st.dataValidation(t, d["Nickname"], "肯·汤普逊")
}

func (st *ServerTest) Update(t *testing.T) {
    // Update
    q := "UPDATE passport_user SET BirthYear = ?, Gender = ?, Nickname = ? WHERE UserID = ?"
    r, err := st.Server.Exec(q, 1982, "Female", "Bob", 1000000)
    if err != nil {
        t.Fatalf("[%s]: %v\n", st.Server.Type, err)
    }
    rowsAffected, err := r.RowsAffected()
    if err != nil {
        t.Fatalf("[%s]: %v\n", st.Server.Type, err)
    }
    if rowsAffected <= 0 {
        t.Fatalf("[%s] Update Failed: %v\n", st.Server.Type, rowsAffected)
    }


    // Update confirm
    d := make(map[string]interface{})
    q = "SELECT * FROM passport_user WHERE UserID = ?"
    err = st.Server.Row(&d, q, 1000000)
    if err != nil {
        t.Fatalf("[%s]: %v\n", st.Server.Type, err)
    }
    st.dataValidation(t, d["CreationTime"], "2015-01-17 00:00:00")
    st.dataValidation(t, d["BirthYear"], int64(1982))
    st.dataValidation(t, d["Gender"], "Female")
    st.dataValidation(t, d["Nickname"], "Bob")
}

func (st *ServerTest) Delete(t *testing.T) {
    // Delete
    q := "DELETE FROM passport_user WHERE UserID = ?"
    r, err := st.Server.Exec(q, 1000000)
    if err != nil {
        t.Fatalf("[%s]: %v\n", st.Server.Type, err)
    }
    rowsAffected, err := r.RowsAffected()
    if err != nil {
        t.Fatalf("[%s]: %v\n", st.Server.Type, err)
    }
    if rowsAffected != 1 {
        t.Fatalf("[%s] Delete Failed: %v\n", st.Server.Type, rowsAffected)
    }
}

func (st *ServerTest) Rows(t *testing.T) {
    d := []map[string]interface{}{}
    q := "SELECT * FROM passport_user WHERE UserID > ?"
    err := st.Server.Rows(&d, q, 1)
    if err != nil {
        t.Fatalf("[%s]: %v\n", st.Server.Type, err)
    }
    if len(d) != 1 {
        t.Fatalf("[%s] Returns the number of rows of data is incorrect: %v\n", st.Server.Type, len(d))
    }
    st.dataValidation(t, d[0]["CreationTime"], "2015-01-17 00:00:00")
    st.dataValidation(t, d[0]["BirthYear"], int64(1982))
    st.dataValidation(t, d[0]["Gender"], "Female")
    st.dataValidation(t, d[0]["Nickname"], "Bob")
}

func (st *ServerTest) dataValidation(t *testing.T, l, r interface{}) {
    if l != r {
        t.Fatalf("[%s] Value validation fails: %v\t%v\n", st.Server.Type, l, r)
    }
}

func NewServerTest(typ string, dsn string) *ServerTest {
    s := ServerTest{}
    s.Server = NewServer(typ, typ, dsn)
    return &s
}

func TestServer(t *testing.T) {
    st := NewServerTest("mysql", "root:@tcp(127.0.0.1:3306)/zhgo?charset=utf8")
    st.Init(t)
    st.Start(t)

    st = NewServerTest("sqlite3", "sqlite3.db")
    st.Init(t)
    st.Start(t)
}
