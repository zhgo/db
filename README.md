# ORM

ORM library for golang. Supported databases include MySQL, MariaDB, PostgreSQL, sqlite3.

[![Build Status](https://travis-ci.org/zhgo/db.svg)](https://travis-ci.org/zhgo/db)
[![Coverage](http://gocover.io/_badge/github.com/zhgo/db)](http://gocover.io/github.com/zhgo/db)
[![GoDoc](https://godoc.org/github.com/zhgo/db?status.png)](http://godoc.org/github.com/zhgo/db)
[![License](https://img.shields.io/badge/license-BSD-ff69b4.svg?style=flat)](https://github.com/zhgo/db/blob/master/LICENSE)

# Overview

## Install

```bash
go get github.com/zhgo/db
```

## Init database

```sql
CREATE TABLE `table1` (
  `UserID` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `CreationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `BirthYear` year(4) NOT NULL,
  `Gender` enum('Secret','Male','Female') NOT NULL DEFAULT 'Secret',
  `Nickname` varchar(16) NOT NULL,
  PRIMARY KEY (`UserID`)
) ENGINE=InnoDB AUTO_INCREMENT=1000000 DEFAULT CHARSET=utf8;
```

## Connect to database

```go
import (
    "github.com/zhgo/db"
)
s := db.NewServer("mysql-1", "mysql", "root:@tcp(127.0.0.1:3306)/zhgo?charset=utf8")
```

## Insert

```go
q := db.NewQuery(s).InsertInto("table1")
r, err := q.Fields("BirthYear", "Gender", "Nickname").Values(1980, "Male", "Bob").Exec()
```

Or:

```go
d = map[string]interface{}{"BirthYear": 1980, "Gender": "Male", "Nickname": "Bob"}
r, err = db.NewQuery(s).InsertInto("table1").Exec(d)
```

## Update

```go
q := db.NewQuery(s)
r, err := q.Update("table1").Set("BirthYear", 1982).Set("Gender", "Female").Set("Nickname", "Bob").Where(q.Eq("UserID", 1000000)).Exec()
```

Or:

```go
d = map[string]interface{}{"BirthYear": 1988, "Gender": "Male", "Nickname": "C语言"}
w = map[string]interface{}{"UserID": 1000001}
r, err = db.NewQuery(s).Update("table1").Exec(d, w)
```

## Delete

```go
q := db.NewQuery(s)
r, err := q.DeleteFrom("table1").Where(q.Eq("UserID", 1000000)).Exec()
```

Or:

```go
w := map[string]interface{}{"UserID": 1000001}
r, err = NewQuery(s).DeleteFrom("table1").Exec(w)
```

## Select

```go
d := []map[string]interface{}{}
q := db.NewQuery(s)
err := q.Select("*").From("table1").Where(q.Eq("UserID", 1000000)).Rows(&d)
```

Or:

```go
d = make(map[string]interface{})
w = map[string]interface{}{"UserID": 1000001}
err = NewQuery(s).Select("*").From("table1").Row(&d, w)
```

# Copyright

Copyright 2015 The zhgo Authors. All rights reserved.

Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.
