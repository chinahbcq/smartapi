/***************************************************************************
 *
 * Copyright (c)2016 SmartApi. All Rights Reserved
 *
 **************************************************************************/


/**
 * @file users.go
 * @author chenqian (chinahbcq@qq.com)
 * @date 2016-05-05 21:37:21
 * @brief 
 *
 **/

package service

import (
	"encoding/json"
	"io"
	"net/http"
	//"time"
)
import (
	"model/data"
	"model/data/dao"
	"utils"
)

type users struct{}

var Users users

func (handle *users) Info(w http.ResponseWriter, r *http.Request, logbuf *utils.LogBuffer) {
	//将query参数打包到map
	m, ok := utils.ParseQuery(r)
	if !ok {
		panic(&utils.SysError{logbuf, "err.param_error"})
	}
	//参数检查
	mustParams := []string{"uid"}
	optParams := []string{"token"}
	ok = utils.CheckParam(logbuf, mustParams, optParams, m)
	if !ok {
		panic(&utils.SysError{logbuf, "err.param_error"})
	}

	uid := m["uid"][0]
	logbuf.WriteLog(" [uid:%s]", uid)
	//参数校验
	Uid := utils.CheckUInt(logbuf, uid)

	var user dao.User
	user, ok = data.User.GetUserById(logbuf, Uid)
	if !ok {
		panic(&utils.SysError{logbuf, "err.user_not_exist"})
	}

	//构造返回内容
	var resp = make(map[string]interface{}, 3)
	resp["request_id"] = logbuf.LogId
	resp["error_code"] = 0

	var profile = make(map[string]interface{}, 3)
	profile["uid"] = user.Uid
	profile["name"] = user.Name
	profile["gender"] = user.Gender

	resp["profile"] = profile

	jsonStr, _ := json.Marshal(resp)
	io.WriteString(w, string(jsonStr))
}
