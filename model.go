// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
    //"strings"
)

// Model struct
type Model struct {
    // Module name, as DB name.
    Module string

    // table instance
    Table Table
}

// Database list
var Servers = make(map[string]*Server)
