/***************************************************************************
 *
 * Copyright (c)2016 SmartApi. All Rights Reserved
 *
 **************************************************************************/


/**
 * @file utils_test.go
 * @author chenqian (chinahbcq@qq.com)
 * @date 2016-04-22 20:04:41
 * @brief utils 单测程序
 *
 **/
package utils

import (
	"encoding/json"
	"testing"
)

var params map[string][]string

func setup() error {
	str := []byte(`{"p1":["val1","val2"], "p2":["go1"]}`)
	err := json.Unmarshal(str, &params)
	if err != nil {
		return err
	}
	return nil
}

func Test_GetOptParam(t *testing.T) {
	err := setup()
	if err != nil {
		t.Error("set up failed")
	}

	t1 := GetOptParam(params, "p3")
	if t1 != "" {
		t.Error("failed")
	}
	t2 := GetOptParam(params, "p1")
	if t2 != "val1" {
		t.Error("failed")
	}
}

func Test_GetSubAction(t *testing.T) {
	var url = "http://abc.com/a/b/c/d.e"
	sub, ok := GetSubAction(url)

	if !ok {
		t.Error("call GetSubAction failed")
	}
	if sub != "e" {
		t.Errorf("the sub action of '%s' is '%s', expected 'e'", url, sub)
	}
}
