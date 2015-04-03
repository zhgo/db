// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
    "errors"
    "fmt"
    "log"
    "reflect"
    "strings"
)

// Query struct
type Query struct {
    // DB
    DB *DB

    // table instance
    Table Table

    // Sql
    sql Sql
}

// Sql struct
type Sql struct {
    //Select fields
    Select string

    //FROM Cation: no use normal
    From string

    //Table join
    Join string

    //Where
    Where string

    //Order by
    Order string

    //Group by
    Group string

    //having
    Having string

    //limit x
    Offset uint64 //limit Offset, Rows

    //limit x, y
    Rows uint64 //limit Offset, Rows

    //lock row(s)
    ForUpdate string

    //args for sql
    Args []interface{}
}

// Select fields.
// Row() or Rows() will return []map[string]interface{} after the method called, otherwith return []EntityStruct.
func (q *Query) Fields(fields ...string) *Query {
    if len(fields) == 0 {
        fields = append(fields, "*")
    }

    q.sql.Select = strings.Join(fields, ", ")
    return q
}

// join table
func (q *Query) Join() *Query {

    return q
}

// Where, 通常查询以Where()方法开始, 后面跟随0个或多个的And(),Or()
func (q *Query) Where(field string, m string, val ...interface{}) *Query {
    return q.whereNode("AND", field, m, val...) //the first where logical must be AND.
}

// AND
func (q *Query) And(field string, m string, val ...interface{}) *Query {
    return q.whereNode("AND", field, m, val...)
}

// OR
func (q *Query) Or(field string, m string, val ...interface{}) *Query {
    return q.whereNode("OR", field, m, val...)
}

// Append where sql, private method
// l logic, for example: AND, OR
// field
// m condition, for example: =, <, <=, >, >=, <>, LLIKE, RLIKE, LIKE, IN, NOT IN, (
// val argument value, 可能是字符串, 也可能是数组, 也可能Where对象
func (q *Query) whereNode(l string, field string, m string, val ...interface{}) *Query {
    switch m {
        case "=", "<", "<=", ">", ">=", "<>", "LIKE":
        q.sql.Where += fmt.Sprintf(" %s %s %s ?", l, field, m)
        q.sql.Args = append(q.sql.Args, val[0])
        case "IN", "NOT IN":
        w := make([]string, len(val))
        for i, v := range val {
            w[i] = "?"
            q.sql.Args = append(q.sql.Args, v)
        }
        q.sql.Where += fmt.Sprintf(" %s %s %s (%s) ", l, field, m, strings.Join(w, ","))
    }

    return q
}

func (q *Query) Order() *Query {
    return q
}

func (q *Query) Group() *Query {

    return q
}

func (q *Query) Having() *Query {

    return q
}

func (q *Query) Limit(offset, rows uint64) *Query {

    return q
}

func (q *Query) Lock() *Query {

    return q
}

// First row, first column
func (q *Query) Scaler(ptr interface{}) error {

    return nil
}

// get one row
func (q *Query) Row(ptr interface{}) error {
    //get first row only
    q.sql.Offset = 0
    q.sql.Rows = 1

    return q.DB.Row(ptr, q.toString(), q.sql.Args)
}

// get all rows, return []struct or []map[string]interface{} or [][]interface{}
func (q *Query) Rows(ptr interface{}) error {
    return q.DB.Rows(ptr, q.toString(), q.sql.Args)
}

// Insert
func (q *Query) Insert(data *[]interface{}) (int64, error) {
    if len(*data) == 0 {
        err := errors.New("empty insert data")
        log.Printf("%s\n", err)
        return 0, err
    }

    var insertVal []string
    insertArgs := []interface{}{}
    insertArgsValue := reflect.Indirect(reflect.ValueOf(&insertArgs))

    for _, item := range *data {
        sli := []string{}

        //convert struct fields to sql string
        itemValue := reflect.ValueOf(item)
        itemType := itemValue.Type()
        for i := 0; i < itemType.NumField(); i++ {
            itemField := itemType.Field(i)
            if !itemField.Anonymous && itemField.Tag.Get("pk") == "" { // && itemField.Tag.Get("json") != ""
                sli = append(sli, "?")
                v := itemValue.FieldByIndex([]int{i})
                insertArgsValue.Set(reflect.Append(insertArgsValue, v))
            }
        }

        insertVal = append(insertVal, "("+strings.Join(sli, ", ")+")")
    }

    insertSql := fmt.Sprintf("INSERT INTO `%s` (`%s`) VALUES %s", q.Table.Name, strings.Join(q.Table.Fields, "`, `"), strings.Join(insertVal, ","))

    log.Printf("%s\n", insertSql)
    log.Printf("%#v\n", insertArgs)

    result, err := q.DB.Exec(insertSql, insertArgs)
    if err != nil {
        log.Printf("%s\n", err)
        return 0, err
    }

    return result.LastInsertId()
}

// Update
func (q *Query) Update(data *map[string]interface{}) (int64, error) {
    //convert map to sql string
    updateBody := []string{}
    updateArgs := []interface{}{}

    for key, element := range *data {
        updateBody = append(updateBody, key+" = ?")
        updateArgs = append(updateArgs, element)
    }

    updateSql := fmt.Sprintf("UPDATE %s SET %s WHERE 1 %s", q.Table.Name, strings.Join(updateBody, ", "), q.sql.Where)

    updateArgs = append(updateArgs, q.sql.Args...)

    log.Printf("%s\n", updateSql)
    log.Printf("%#v\n", updateArgs)

    result, err := q.DB.Exec(updateSql, updateArgs)
    if err != nil {
        log.Printf("%s\n", err)
        return 0, err
    }

    return result.RowsAffected()
}

// Delete
func (q *Query) Delete() (int64, error) {
    deleteSql := fmt.Sprintf("DELETE FROM %s WHERE 1 %s", q.Table.Name, q.sql.Where)

    log.Printf("%s\n", deleteSql)
    log.Printf("%#v\n", q.sql.Args)

    result, err := q.DB.Exec(deleteSql, q.sql.Args)
    if err != nil {
        log.Printf("%s\n", err)
        return 0, err
    }

    return result.RowsAffected()
}

// connect all sql department to a corect sql string.
func (q *Query) toString() string {
    //create sql
    var sel string
    if len(q.sql.Select) > 0 {
        sel = q.sql.Select //返回[]map[string]interface{}
    } else {
        sel = fmt.Sprintf("a.`%s`, a.`%s`", q.Table.Primary, strings.Join(q.Table.Fields, "`, a.`")) //返回[]EntityStruct
    }

    sql := fmt.Sprintf("SELECT %s FROM %s a %s WHERE 1 %s %s %s %s LIMIT %d, %d %s", sel, q.Table.Name, q.sql.Join, q.sql.Where, q.sql.Order, q.sql.Group, q.sql.Having, q.sql.Offset, q.sql.Rows, q.sql.ForUpdate)

    log.Printf("%s\n", sql)
    log.Printf("%#v\n", q.sql.Args)

    return sql
}
