/***************************************************************************
 *
 * Copyright (c)2016 SmartApi. All Rights Reserved
 *
 **************************************************************************/


/**
 * @file log.go
 * @author chenqian (chinahbcq@qq.com)
 * @date 2016-04-30 20:25:05
 * @brief 
 *
 **/

package utils

import (
	"bytes"
	"fmt"
)

type LogBuffer struct {
	buf   bytes.Buffer
	LogId int64
}

func (handle *LogBuffer) WriteLog(format string, a ...interface{}) (err error) {
	handle.buf.WriteString(fmt.Sprintf(format, a...))
	return nil
}
func (handle *LogBuffer) String() string {
	return handle.buf.String()
}
func NewLogBuffer() *LogBuffer {
	var buf bytes.Buffer
	lb := LogBuffer{buf, GenLogId()}
	lb.WriteLog("[logid:%d]", lb.LogId)
	return &lb
}
