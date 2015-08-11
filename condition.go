// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import ()

// Condition struct
type Condition struct {
	Eq   map[string]string
	Ge   map[string]string
	Gt   map[string]string
	Le   map[string]string
	Lt   map[string]string
	Ne   map[string]string
	Like map[string]string
	In   map[string][]string
}

func (c *Condition) Success() {

}

// New condition
func NewCondition() {

}
