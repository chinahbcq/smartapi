/***************************************************************************
 *
 * Copyright (c)2016 SmartApi. All Rights Reserved
 *
 **************************************************************************/


/**
 * @file log_test.go
 * @author chenqian (chinahbcq@qq.com)
 * @date 2016-04-30 20:35:16
 * @brief 
 *
 **/
package utils

import (
	"fmt"
	"testing"
)

func revoke(logs *LogBuffer) {
	logs.WriteLog("this is %s, and %d\n", "test", 8)
}
func Test_Log(t *testing.T) {
	var logs = NewLogBuffer()

	logs.WriteLog("this is %s, and %d", "test", 9)
	logs.WriteLog("this is %s, and %d", "test", 9)
	revoke(logs)
	fmt.Println(logs.String())
}
