# ORM library

Supported databases include MySQL, MariaDB, PostgreSQL, Sqlite3.

[![Build Status](https://travis-ci.org/zhgo/db.svg)](https://travis-ci.org/zhgo/db)
[![Coverage Status](https://coveralls.io/repos/zhgo/db/badge.svg)](https://coveralls.io/r/zhgo/db)
[![GoDoc](https://godoc.org/github.com/zhgo/db?status.png)](http://godoc.org/github.com/zhgo/db)
[![License](https://img.shields.io/badge/license-BSD-blue.svg?style=flat)](https://github.com/zhgo/db/blob/master/LICENSE)

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

Or:

```go
s := db.Connect("mysql-1", "mysql", "root:@tcp(127.0.0.1:3306)/zhgo?charset=utf8")
```

**mysql-1** is connection name. **mysql** is sql driver type. **root:@tcp(127.0.0.1:3306)/zhgo?charset=utf8** is DSN.

## Insert

```go
// INSERT INTO table1(BirthYear, Gender, Nickname) VALUES(1980, 'Male', 'Bob')
q := s.NewQuery()
r, err := q.InsertInto("table1").Fields("BirthYear", "Gender", "Nickname").Values(1980, "Male", "Bob").Exec()
```

**Values()** method can be called multiple times to insert multiple rows.

Or:

```go
// INSERT INTO table1(BirthYear, Gender, Nickname) VALUES(1980, 'Male', 'Bob')
d = db.Item{"BirthYear": 1980, "Gender": "Male", "Nickname": "Bob"}
r, err := s.InsertInto("table1").Exec(d)
```

**d** is a map type.

## Update

```go
// UPDATE table1 SET BirthYear = 1982, Gender = 'Female', Nickname = 'Bob' WHERE UserID = 1000000
q := s.NewQuery()
q.Update("table1")
q.Set("BirthYear", 1982)
q.Set("Gender", "Female")
q.Set("Nickname", "Bob")
r, err := q.Where(q.Eq("UserID", 1000000)).Exec()
```

**Set()** method can be called multiple times to update multiple fields.

Or:

```go
// UPDATE table1 SET BirthYear = 1988, Gender = 'Male', Nickname = 'C语言' WHERE UserID = 1000001
d = db.Item{"BirthYear": 1988, "Gender": "Male", "Nickname": "C语言"}
w = db.Where{"UserID": 1000001}
r, err := s.Update("table1").Exec(d, w)
```

## Delete

```go
// DELETE FROM table1 WHERE UserID = 1000000
q := s.NewQuery()
r, err := q.DeleteFrom("table1").Where(q.Eq("UserID", 1000000)).Exec()
```

Or:

```go
// DELETE FROM table1 WHERE UserID = 1000001
w := db.Where{"UserID": 1000001}
r, err := s.DeleteFrom("table1").Exec(w)
```

## Select

```go
// SELECT * FROM table1 WHERE UserID = 1000001
d := []db.Item{}
q := s.NewQuery()
err := q.Select("*").From("table1").Where(q.Eq("UserID", 1000000)).Rows(&d)
```

Or:

```go
// SELECT * FROM table1 WHERE UserID = 1000001
d = make(db.Item)
w = db.Where{"UserID": 1000001}
err := s.Select("*").From("table1").Row(&d, w)
```

# Copyright

Copyright 2015 The zhgo Authors. All rights reserved.

Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.
