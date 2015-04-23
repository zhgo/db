# ORM library for golang

Supported databases include MySQL, MariaDB, PostgreSQL, Sqlite3.

[![Build Status](https://travis-ci.org/zhgo/db.svg)](https://travis-ci.org/zhgo/db)
[![Coverage](http://gocover.io/_badge/github.com/zhgo/db)](http://gocover.io/github.com/zhgo/db)
[![GoDoc](https://godoc.org/github.com/zhgo/db?status.png)](http://godoc.org/github.com/zhgo/db)
[![License](https://img.shields.io/badge/license-BSD-ff69b4.svg?style=flat)](https://github.com/zhgo/db/blob/master/LICENSE)

# Useage

## Install

```bash
go get github.com/zhgo/db
```

## Sample sql scripts

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

The example use MySQL as default database instance.

## Import

```go
import (
    "github.com/zhgo/db"
)
```

## Connect to database

```go
s := db.NewServer("mysql-1", "mysql", "root:@tcp(127.0.0.1:3306)/zhgo?charset=utf8")
```

**mysql-1** is connection name. **mysql** is sql driver type. **root:@tcp(127.0.0.1:3306)/zhgo?charset=utf8** is DSN.

## Insert

```go
q := db.NewQuery(s).InsertInto("table1")
r, err := q.Fields("BirthYear", "Gender", "Nickname").Values(1980, "Male", "Bob").Exec()
```

**Values()** method can be called multiple times to insert multiple rows.

Or:

```go
d = db.Item{"BirthYear": 1980, "Gender": "Male", "Nickname": "Bob"}
r, err = db.NewQuery(s).InsertInto("table1").Exec(d)
```

**d** is a map type.

## Update

```go
q := db.NewQuery(s).Update("table1")
r, err := q.Set("BirthYear", 1982).Set("Gender", "Female").Set("Nickname", "Bob").Where(q.Eq("UserID", 1000000)).Exec()
```

**Set()** method can be called multiple times to update multiple fields.

Or:

```go
d = db.Item{"BirthYear": 1988, "Gender": "Male", "Nickname": "C语言"}
w = db.Where{"UserID": 1000001}
r, err = db.NewQuery(s).Update("table1").Exec(d, w)
```

## Delete

```go
q := db.NewQuery(s).DeleteFrom("table1")
r, err := q.Where(q.Eq("UserID", 1000000)).Exec()
```

Or:

```go
w := db.Where{"UserID": 1000001}
r, err = db.NewQuery(s).DeleteFrom("table1").Exec(w)
```

## Select

```go
d := []db.Item{}
q := db.NewQuery(s).Select("*")
err := q.From("table1").Where(q.Eq("UserID", 1000000)).Rows(&d)
```

Or:

```go
d = make(db.Item)
w = db.Where{"UserID": 1000001}
err = db.NewQuery(s).Select("*").From("table1").Row(&d, w)
```

# Copyright

Copyright 2015 The zhgo Authors. All rights reserved.

Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.
