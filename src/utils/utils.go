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
	"time"
)

type SysError struct {
	LogBuf  *LogBuffer
	ErrInfo string
}

type SysErrorExt struct {
	LogBuf  *LogBuffer
	ErrInfo string
	More    map[string]interface{}
}

func (e *SysErrorExt) Error() string {
	info := ErrorMap[e.ErrInfo]
	v := info.(map[string]interface{})
	v["request_id"] = e.LogBuf.LogId

	for key, val := range e.More {
		v[key] = val
	}
	Global.Logger.UbLogNotice(e.LogBuf.String())
	str, _ := json.Marshal(v)
	return string(str)
}

func (e *SysError) Error() string {
	info := ErrorMap[e.ErrInfo]
	v := info.(map[string]interface{})
	v["request_id"] = e.LogBuf.LogId

	Global.Logger.UbLogNotice(e.LogBuf.String())
	str, _ := json.Marshal(v)
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

func CheckParam(mustParams, optParams []string, m map[string][]string) bool {
	for _, param := range mustParams {
		if _, ok := m[param]; !ok || len(m[param][0]) < 1 {
			return false
		}
	}

	for _, param := range optParams {
		if _, ok := m[param]; ok {
			if len(m[param][0]) < 1 {
				return false
			}
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
