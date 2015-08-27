// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
    "fmt"
    "io/ioutil"
    "strings"
    "testing"
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
    qt.Delete(t)
}

func (qt *QueryTest) Insert(t *testing.T) {
    // Insert
    q := qt.Query.Server.InsertInto("passport_user")
    q.SetPrimary("UserID") // PostgreSQL compatibility
    r, err := q.Fields("UserID", "CreationTime", "BirthYear", "Gender", "Nickname").Values(1000000, "2015-01-17 00:00:00", 1980, "Male", "肯·汤普逊").Exec()
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    if r.LastInsertId != 1000000 {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, r.LastInsertId)
    }

    // Insert confirm
    d := make(Item)
    q = NewQuery(qt.Query.Server)
    err = q.Select("*").From("passport_user").Where(q.Eq("UserID", 1000000)).Row(&d)
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    qt.dataValidation(t, string(d["CreationTime"].([]byte)), "2015-01-17 00:00:00")
    qt.dataValidation(t, d["BirthYear"], int64(1980))
    qt.dataValidation(t, string(d["Gender"].([]byte)), "Male")
    qt.dataValidation(t, string(d["Nickname"].([]byte)), "肯·汤普逊")

    // Insert
    d = Item{
        "UserID":       1000001,
        "CreationTime": "2015-01-17 01:00:00",
        "BirthYear":    1986,
        "Gender":       "Secret",
        "Nickname":     "阿里马马"}
    q = NewQuery(qt.Query.Server)
    q.SetPrimary("UserID") // PostgreSQL compatibility
    r, err = q.InsertInto("passport_user").Exec(d)
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    if r.LastInsertId != 1000001 {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, r.LastInsertId)
    }

    // Insert confirm
    d = make(Item)
    q = NewQuery(qt.Query.Server)
    err = q.Select("*").From("passport_user").Where(q.Eq("UserID", 1000001)).Row(&d)
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    qt.dataValidation(t, string(d["CreationTime"].([]byte)), "2015-01-17 01:00:00")
    qt.dataValidation(t, d["BirthYear"], int64(1986))
    qt.dataValidation(t, string(d["Gender"].([]byte)), "Secret")
    qt.dataValidation(t, string(d["Nickname"].([]byte)), "阿里马马")
}

func (qt *QueryTest) Update(t *testing.T) {
    // Update
    q := NewQuery(qt.Query.Server)
    r, err := q.Update("passport_user").Set("BirthYear", 1982).Set("Gender", "Female").Set("Nickname", "Bob").Where(q.Eq("UserID", 1000000)).Exec()
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    if r.RowsAffected != 1 {
        t.Fatalf("[%s] Update Failed: %v\n", qt.Query.Server.Type, r.RowsAffected)
    }

    // Update confirm
    d := make(Item)
    w := Where{"UserID": 1000000}
    err = qt.Query.Server.Select("*").From("passport_user").Row(&d, w)
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    qt.dataValidation(t, string(d["CreationTime"].([]byte)), "2015-01-17 00:00:00")
    qt.dataValidation(t, d["BirthYear"], int64(1982))
    qt.dataValidation(t, string(d["Gender"].([]byte)), "Female")
    qt.dataValidation(t, string(d["Nickname"].([]byte)), "Bob")

    // Update
    d = Item{
        "BirthYear": 1988,
        "Gender":    "Male",
        "Nickname":  "C语言"}
    w = Where{"UserID": 1000001}
    r, err = qt.Query.Server.Update("passport_user").Exec(d, w)
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    if r.RowsAffected != 1 {
        t.Fatalf("[%s] Update Failed: %v\n", qt.Query.Server.Type, r.RowsAffected)
    }

    // Update confirm
    d = make(Item)
    q = NewQuery(qt.Query.Server)
    err = q.Select().From("passport_user").Row(&d, w)
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    qt.dataValidation(t, string(d["CreationTime"].([]byte)), "2015-01-17 01:00:00")
    qt.dataValidation(t, d["BirthYear"], int64(1988))
    qt.dataValidation(t, string(d["Gender"].([]byte)), "Male")
    qt.dataValidation(t, string(d["Nickname"].([]byte)), "C语言")
}

func (qt *QueryTest) Delete(t *testing.T) {
    // Delete
    q := NewQuery(qt.Query.Server)
    r, err := q.DeleteFrom("passport_user").Where(q.Eq("UserID", 1000000)).Exec()
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    if r.RowsAffected != 1 {
        t.Fatalf("[%s] Delete Failed: %v\n", qt.Query.Server.Type, r.RowsAffected)
    }

    // Delete
    w := Where{"UserID": 1000001}
    r, err = qt.Query.Server.DeleteFrom("passport_user").Exec(w)
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    if r.RowsAffected != 1 {
        t.Fatalf("[%s] Delete Failed: %v\n", qt.Query.Server.Type, r.RowsAffected)
    }
}

func (qt *QueryTest) Rows(t *testing.T) {
    // Rows
    d := []Item{}
    q := NewQuery(qt.Query.Server)
    err := q.Select("*").From("passport_user").Where(q.In("UserID", 1000000, 1000001)).Rows(&d)
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    if len(d) != 2 {
        t.Fatalf("[%s] Returns the number of rows of data is incorrect: %v\n", qt.Query.Server.Type, len(d))
    }
    qt.dataValidation(t, string(d[0]["CreationTime"].([]byte)), "2015-01-17 00:00:00")
    qt.dataValidation(t, d[0]["BirthYear"], int64(1982))
    qt.dataValidation(t, string(d[0]["Gender"].([]byte)), "Female")
    qt.dataValidation(t, string(d[0]["Nickname"].([]byte)), "Bob")

    // Rows
    d = []Item{}
    q = NewQuery(qt.Query.Server)
    err = q.Select().From("passport_user").Where(q.Eq("UserID", 1000000), q.Or(q.Eq("UserID", 1000001))).Rows(&d)
    if err != nil {
        t.Fatalf("[%s]: %v\n", qt.Query.Server.Type, err)
    }
    if len(d) != 2 {
        t.Fatalf("[%s] Returns the number of rows of data is incorrect: %v\n", qt.Query.Server.Type, len(d))
    }
    qt.dataValidation(t, string(d[0]["CreationTime"].([]byte)), "2015-01-17 00:00:00")
    qt.dataValidation(t, d[0]["BirthYear"], int64(1982))
    qt.dataValidation(t, string(d[0]["Gender"].([]byte)), "Female")
    qt.dataValidation(t, string(d[0]["Nickname"].([]byte)), "Bob")
}

func (qt *QueryTest) dataValidation(t *testing.T, l, r interface{}) {
    if l != r {
        t.Fatalf("[%s] Value validation fails: %v\t%v\n", qt.Query.Server.Type, l, r)
    }
}

func NewQueryTest(typ string, dsn string) *QueryTest {
    s := NewServer(typ, dsn)
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

    st = NewQueryTest("postgres", "user=postgres dbname=zhgo sslmode=disable")
    //st = NewServerTest("postgres", "postgres://LD:@localhost:5432/zhgo?sslmode=verify-full")
    st.Init(t)
    st.Load(t)
    st.Start(t)
}
