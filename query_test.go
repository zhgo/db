// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
    "testing"
    "io/ioutil"
    "fmt"
    "strings"
)

type QueryTest struct {
    // Query
    Query *Query
}

func (qt *QueryTest) Init(t *testing.T) {
    // Load sql file
    p := fmt.Sprintf("tb-%s.sql", qt.Query.Server.Type)
    b, err := ioutil.ReadFile(p)
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    initSQL := strings.Split(string(b), ";")

    // Run Init SQL
    for _, v := range initSQL {
        if strings.Trim(v, "\r\n \t") != "" {
            _, err := qt.Query.Server.Exec(v)
            if err != nil {
                t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
            }
        }
    }
}

func (qt *QueryTest) Load(t *testing.T) {

}

func (qt *QueryTest) Start(t *testing.T) {
    qt.Insert(t)
    qt.Update(t)
    qt.Rows(t)
}

func (qt *QueryTest) Insert(t *testing.T) {
    // Insert
    q := NewQuery(qt.Query.Server)
    q.InsertInto("passport_user")
    q.Fields("UserID", "CreationTime", "BirthYear", "Gender", "Nickname")
    q.Values(1000000, "2015-01-17 00:00:00", 1980, "Male", "肯·汤普逊")
    r, err := q.Exec()
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    lastInsertId, err := r.LastInsertId()
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    if lastInsertId != 1000000 {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, lastInsertId)
    }


    // Insert confirm
    d := make(map[string]interface{})
    q = NewQuery(qt.Query.Server)
    err = q.Select("*").From("passport_user").Where(q.Eq("UserID", 1000000)).Row(&d)
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    qt.dataValidation(t, d["CreationTime"], "2015-01-17 00:00:00")
    qt.dataValidation(t, d["BirthYear"], int64(1980))
    qt.dataValidation(t, d["Gender"], "Male")
    qt.dataValidation(t, d["Nickname"], "肯·汤普逊")
}

func (qt *QueryTest) Update(t *testing.T) {
    // Update
    q := NewQuery(qt.Query.Server)
    r, err := q.Update("passport_user").Set("BirthYear", 1982).Set("Gender", "Female").Set("Nickname", "Bob").Where(q.Eq("UserID", 1000000)).Exec()
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    rowsAffected, err := r.RowsAffected()
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    if rowsAffected <= 0 {
        t.Fatalf("[%s] Update Failed: %v\n", qt.Query.Server.Type, rowsAffected)
    }


    // Update confirm
    d := make(map[string]interface{})
    q = NewQuery(qt.Query.Server)
    err = q.Select("*").From("passport_user").Where(q.Eq("UserID", 1000000)).Row(&d)
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    qt.dataValidation(t, d["CreationTime"], "2015-01-17 00:00:00")
    qt.dataValidation(t, d["BirthYear"], int64(1982))
    qt.dataValidation(t, d["Gender"], "Female")
    qt.dataValidation(t, d["Nickname"], "Bob")
}

func (qt *QueryTest) Rows(t *testing.T) {
    d := []map[string]interface{}{}
    q := NewQuery(qt.Query.Server)
    err := q.Select("*").From("passport_user").Where(q.Gt("UserID", 1)).Rows(&d)
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    if len(d) != 1 {
        t.Fatalf("[%s] Returns the number of rows of data is incorrect: %v\n", qt.Query.Server.Type, len(d))
    }
    qt.dataValidation(t, d[0]["CreationTime"], "2015-01-17 00:00:00")
    qt.dataValidation(t, d[0]["BirthYear"], int64(1982))
    qt.dataValidation(t, d[0]["Gender"], "Female")
    qt.dataValidation(t, d[0]["Nickname"], "Bob")
}

func (qt *QueryTest) dataValidation(t *testing.T, l, r interface{}) {
    if l != r {
        t.Fatalf("[%s] Value validation fails: %v\t%v\n", qt.Query.Server.Type, l, r)
    }
}

func NewQueryTest(typ string, dsn string) *QueryTest {
    s := NewServer(typ, typ, dsn)
    qt := QueryTest{}
    qt.Query = NewQuery(s)
    return &qt
}

func TestQuery(t *testing.T) {
    st := NewQueryTest("mysql", "root:@tcp(127.0.0.1:3306)/zhgo?charset=utf8")
    st.Init(t)
    st.Load(t)
    st.Start(t)

    st = NewQueryTest("sqlite3", "sqlite3.db")
    st.Init(t)
    st.Load(t)
    st.Start(t)
}
