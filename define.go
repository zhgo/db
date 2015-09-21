// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
    "database/sql"
)

// Alias of map[string]interface{}
type Item map[string]interface{}

// Alias of []map[string]interface{}
type Items []map[string]interface{}

// Alias of map[string]interface{}
type Where map[string]interface{}

// Enviroment: 0, 1, 2, 3
var Env int8 = 0

// Server instance
var dbObjects map[string]*sql.DB = make(map[string]*sql.DB)

// Server list
var Servers = make(map[string]*Server)
