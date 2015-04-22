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
s := NewServer("mysql-1", "mysql", "root:@tcp(127.0.0.1:3306)/zhgo?charset=utf8")
```

## Insert example

```go
q := NewQuery(s).InsertInto("table1")
r, err := q.Fields("BirthYear", "Gender", "Nickname").Values(1980, "Male", "Bob").Exec()
```

Or:

```go
d = map[string]interface{}{"BirthYear": 1980, "Gender": "Male", "Nickname": "Bob"}
r, err = NewQuery(s).InsertInto("table1").Exec(d)
```

## Select example

```go
d := []map[string]interface{}{}
q := NewQuery(s)
err := q.Select("*").From("table1").Where(q.Eq("UserID", 1000000)).Rows(&d)
```

