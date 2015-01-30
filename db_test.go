// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"reflect"
	"testing"
)

type TestRegion struct {
	RegionId       int64 `pk:"true"`
	ParentRegionId int64
	Title          string
}

type TestRegionTag struct {
	RegionId       int64  `json:"region_id" pk:"true"`
	ParentRegionId int64  `json:"parent_region_id"`
	Title          string `json:"title"`
}

func TestDataGetScanVariable(t *testing.T) {
	columnsLen := 3

	structSli := []TestRegion{}
	elemKind, scanPtr, scanArgs, err := scanVariable(&structSli, columnsLen, true)
	if err != nil {
		t.Errorf("GetScanVariable error: %v", err)
	}
	if elemKind != reflect.Struct {
		t.Errorf("GetScanVariable elemKind error: %v", elemKind)
	}
	if scanPtr == nil {
		t.Error("GetScanVariable scanPtr error")
	}
	if scanArgs == nil {
		t.Error("GetScanVariable scanArgs error")
	}
}

func TestDataTypeAssertion(t *testing.T) {
	//int8
	/*var srcInt8 int8 = 8
	if dstInt8 := typeAssertion(srcInt8); dstInt8 != srcInt8 {
		t.Error("Int8 type assert error")
	}*/

	//int64
	var srcInt64 int64 = 128
	if dstInt64 := typeAssertion(srcInt64); dstInt64 != srcInt64 {
		t.Error("Int64 type assert error")
	}

	//[]string
	/*srcSlice := []string{"a"}
	dstSlice := typeAssertion(srcSlice)
	if dstSlice != srcSlice {
		t.Error("Slice type assert error")
	}*/
}

func TestDataGetFields(t *testing.T) {
	primary, fields := TableFields(new(TestRegion))
	if primary != "RegionId" {
		t.Errorf("primary detected failure: %#v", primary)
	}
	if len(fields) != 2 {
		t.Error("fields detected failure")
	}
	if fields[0] != "ParentRegionId" || fields[1] != "Title" {
		t.Errorf("fields detected failure: %#v", fields)
	}

	primary, fields = TableFields(new(TestRegionTag))
	if primary != "region_id" {
		t.Errorf("primary detected failure: %#v", primary)
	}
	if len(fields) != 2 {
		t.Error("fields detected failure")
	}
	if fields[0] != "parent_region_id" || fields[1] != "title" {
		t.Errorf("fields detected failure: %#v", fields)
	}
}
