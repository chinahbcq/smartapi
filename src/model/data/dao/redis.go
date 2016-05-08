/***************************************************************************
 *
 * Copyright (c)2016 SmartApi. All Rights Reserved
 *
 **************************************************************************/


/**
 * @file dao/redis.go
 * @author chenqian (chinahbcq@qq.com)
 * @date 2016-05-03 21:36:48
 * @brief 
 *
 **/
package dao

import (
	"gopkg.in/redis.v3"
	"time"
	"utils"
)

type redisHandle struct{}

var RedisHandle redisHandle

func (handle *redisHandle) Get(logbuf *utils.LogBuffer, key string) (val string, err error) {
	startTime := time.Now()
	logger := utils.Global.Logger
	redisCli, _ := utils.Global.RedisHandles.GetClient()
	val, err = redisCli.Get(key).Result()
	if err == redis.Nil {
		logbuf.WriteLog(" [get_redis:redis.Nil]")
		logbuf.WriteLog(" [get_redis_cost_time:%v]", time.Now().Sub(startTime))
		return
	} else if err != nil {
		logbuf.WriteLog(" [get_redis:ERROR]")
		logger.UbLogWarning("cache_bduss:REDIS_SERVER_ERR, msg:%s", err.Error())
	} else {
		logbuf.WriteLog(" [get_redis:OK]")
		err = nil
	}
	logbuf.WriteLog(" [get_redis_cost_time:%v]", time.Now().Sub(startTime))
	return
}

func (handle *redisHandle) Set(logbuf *utils.LogBuffer, key, val string, expire_s int) error {
	startTime := time.Now()
	redisCli, _ := utils.Global.RedisHandles.GetClient()

	duration := time.Duration(expire_s) * time.Second
	err := redisCli.Set(key, val, duration).Err()
	if err != nil {
		logbuf.WriteLog(" [set_redis_fail:%s]", err.Error())
		logbuf.WriteLog(" [set_redis_cost_time:%v]", time.Now().Sub(startTime))
		return err
	}
	logbuf.WriteLog(" [set_redis:OK]")
	logbuf.WriteLog(" [set_redis_cost_time:%v]", time.Now().Sub(startTime))
	return nil
}

func (handle *redisHandle) HSet(logbuf *utils.LogBuffer, key, field, val string) error {
	startTime := time.Now()
	redisCli, _ := utils.Global.RedisHandles.GetClient()

	err := redisCli.HSet(key, field, val).Err()
	if err != nil {
		logbuf.WriteLog(" [hset_redis_fail:%s]", err.Error())
		logbuf.WriteLog(" [hset_redis_cost_time:%v]", time.Now().Sub(startTime))
		return err
	}
	logbuf.WriteLog(" [hset_redis:OK]")
	logbuf.WriteLog(" [hset_redis_cost_time:%v]", time.Now().Sub(startTime))
	return nil
}
