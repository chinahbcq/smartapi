/***************************************************************************
 *
 * Copyright (c)2016 SmartApi. All Rights Reserved
 *
 **************************************************************************/


/**
 * @file user.go
 * @author chenqian (chinahbcq@qq.com)
 * @date 2016-05-07 21:21:52
 * @brief 
 *
 **/

package data

import (
	"model/data/dao"
	"utils"
)

type user struct{}

var User user

func (handle *user) GetUserById(logbuf *utils.LogBuffer, uid int64) (dao.User, bool) {
	return dao.MysqlHandle.GetUser(logbuf, uid)
}
