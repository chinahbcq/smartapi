/***************************************************************************
 *
 * Copyright (c)2016 SmartApi. All Rights Reserved
 *
 **************************************************************************/

/**
 * @file utils.go
 * @author chenqian (chinahbcq@qq.com)
 * @date 2016-04-19 14:32:07
 * @brief
 *
 **/

package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type SysError struct {
	LogBuf  *LogBuffer
	ErrInfo string
}

func (e *SysError) Error() string {
	info := ErrorMap[e.ErrInfo]
	info.RequestId = e.LogBuf.LogId
	e.LogBuf.WriteLog(" [error_msg:%s] [error_code:%d]", e.ErrInfo, info.ErrorCode)
	Global.Logger.UbLogNotice(e.LogBuf.String())

	str, _ := json.Marshal(info)
	return string(str)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetSubAction(url string) (string, bool) {
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] == '.' {
			subAction := url[i+1:]
			return subAction, true
		}
	}
	return "", false
}

func GenLogId() int64 {
	return time.Now().UnixNano() / 1000000
}

func ParseQuery(r *http.Request) (m map[string][]string, ok bool) {
	ok = false
	if r.Method == "GET" {
		var u = r.URL
		var err error
		if m, err = url.ParseQuery(u.RawQuery); err != nil {
			ok = false
			return
		}
		ok = true
	} else if r.Method == "POST" {
		r.ParseForm()
		m = r.Form
		ok = true
	}
	return
}

func CheckParam(logbuf *LogBuffer, mustParams, optParams []string, m map[string][]string) bool {
	for _, param := range mustParams {
		logbuf.WriteLog(" [%s:", param)
		if _, ok := m[param]; !ok || len(m[param][0]) < 1 {
			logbuf.WriteLog("%s]", "")
			return false
		}
		logbuf.WriteLog("%s]", m[param][0])
	}

	for _, param := range optParams {
		logbuf.WriteLog(" [%s:", param)
		if _, ok := m[param]; ok {
			logbuf.WriteLog("%s]", m[param][0])
			if len(m[param][0]) < 1 {
				return false
			}
		} else {
			logbuf.WriteLog("%s]", "")
		}
	}

	return true
}

func GetOptParam(m map[string][]string, key string) string {
	if _, ok := m[key]; !ok || len(m[key][0]) < 1 {
		return ""
	}
	return m[key][0]
}

func CheckInt(logbuf *LogBuffer, param string) int64 {
	num, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		logbuf.WriteLog(" [error_msg:param '%s' not numeric]", param)
		panic(&SysError{logbuf, "err.param_not_num"})
	}
	return num
}

func CheckUInt(logbuf *LogBuffer, param string) int64 {
	num, err := strconv.ParseUint(param, 10, 64)
	if err != nil || num < 1 {
		logbuf.WriteLog(" [error_msg:param '%s' should uint64 and > 0]", param)
		panic(&SysError{logbuf, "err.param_not_uint"})
	}
	return int64(num)
}
