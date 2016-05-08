/***************************************************************************
 *
 * Copyright (c)2016 SmartApi. All Rights Reserved
 *
 **************************************************************************/


/**
 * @file dao/mysql.go
 * @author chenqian (chinahbcq@qq.com)
 * @date 2016-05-07 21:15:01
 * @brief 
 *
 **/

package dao

import (
	"database/sql"
	"utils"
)

type mysqlHandle struct{}

var MysqlHandle mysqlHandle

func (handle *mysqlHandle) GetUser(logbuf *utils.LogBuffer, uid int64) (user User, ok bool) {
	client, _ := utils.Global.MysqlHandles.GetClient()
	logger := utils.Global.Logger
	err := client.QueryRow("select uid,name,gender from user_info where uid = ?", uid).Scan(&user.Uid,
		&user.Name,
		&user.Gender)

	ok = true
	if err == sql.ErrNoRows {
		logger.UbLogWarning("%s", err.Error())
		ok = false
		return
	} else if err != nil {
		logger.UbLogWarning("%s", err.Error())
		ok = false
		return
	}

	return
}
