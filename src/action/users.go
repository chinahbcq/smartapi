/***************************************************************************
 *
 * Copyright (c)2016 SmartApi. All Rights Reserved
 *
 **************************************************************************/


/**
 * @file action/users.go
 * @author chenqian (chinahbcq@qq.com)
 * @date 2016-05-05 20:31:43
 * @brief 
 *
 **/

package action

import (
	"model/service"
	"net/http"
	"utils"
)

func Users(w http.ResponseWriter, r *http.Request) {
	logger := utils.Global.Logger
	logbuf := utils.NewLogBuffer()
	subAction, ok := utils.GetSubAction(r.URL.Path)
	if !ok {
		panic(&utils.SysError{logbuf, "err.method_not_support"})
	}

	switch subAction {
	case "info":
		service.Users.Info(w, r, logbuf)
	default:
		panic(&utils.SysError{logbuf, "err.method_not_support"})
	}

	logbuf.WriteLog(" [error_code:%d]", 0)
	logger.UbLogNotice("%s", logbuf.String())
}
