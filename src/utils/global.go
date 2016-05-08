/***************************************************************************
 *
 * Copyright (c)2016 SmartApi. All Rights Reserved
 *
 **************************************************************************/


/**
 * @file global.go
 * @author chenqian (chinahbcq@qq.com)
 * @date 2016-04-19 14:28:03
 * @brief 
 *
 **/

package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ublog"
	"gopkg.in/redis.v3"
	"math/rand"
	"time"
)

//定义mysql，redis的handle
type RedisHandle []*redis.Client
type MysqlHandle []*sql.DB

type GlobalStruct struct {
	Config       ConfigStruct
	Logger       *ublog.UbLog
	MysqlHandles MysqlHandle
	RedisHandles RedisHandle
}

type ConfigStruct struct {
	//talking options
	ListenAddress string
	ProcessNum    int

	//log options
	LogDir            string `flag:"log-path"`
	LogFile           string
	LogLevel          uint32 `flag:"log-level"`
	LogChannelLen     uint32 `flag:"log-channel-len"`
	LogFlushThreadNum uint32 `flag:"log-flush-threadnum"`
	LogFlushInterval  uint32 `flag:"log-flush-interval"`

	//net options
	ReadTimeout  int `flag:"read-timeout"`
	WriteTimeout int `flag:"write-timeout"`
	IdleTimeout  int
	RetryTimes   int

	//mongodb options
	ImMongoDBAddresses  []string `flag:"immongodb-addresses"`
	MsgMongoDBAddresses []string `flag:"msgmongodb-addresses"`
	MongoDBConnTimeout  int      `flag:"mongodb-conn-timeout"`
	MongoDBOpTimeout    int      `flag:"mongodb-op-timeout"`
	MongoDBPoolLimit    int      `flag:"mongodb-pool-limit"`
	MongoDBBatchSize    int      `flag:"mongodb-batch-size"`
	MongoDBRetryTimes   int      `flag:"mongodb-retry-times"`

	//mysql options
	MysqlUsername  string
	MysqlPassword  string
	MysqlAddresses []string
	MysqlDB        string

	//redis options
	RedisDialTimeout  int      `flag:"redis-conn-timeout"`
	RedisReadTimeout  int      `flag:"redis-read-timeout"`
	RedisWriteTimeout int      `flag:"redis-write-timeout"`
	RedisPoolSize     int      `flag:"redis-pool-size"`
	RedisAddresses    []string `flag:"redis-addresses"`
}

var Global GlobalStruct
var ErrorMap map[string]interface{}

/**
 * @brief 获取redis的一个实例
 * @param
 * @return
 * @author chinahbcq@qq.com
 * @date 2016-04-22 17:06:00
 */
func (handle *RedisHandle) GetClient() (*redis.Client, bool) {
	idx := rand.Intn(len(*handle))
	return (*handle)[idx], true
}

/**
 * @brief 给redis handle增加一个client实例
 * @param
 * @return
 * @author chinahbcq@qq.com
 * @date 2016-04-22 16:52:04
 */
func (handle *RedisHandle) Append(client *redis.Client) {
	*handle = append(*handle, client)
}

/**
 * @brief 初始化redis handle
 * @param
 * @return
 * @author chinahbcq@qq.com
 * @date 2016-04-22 14:28:57
 */
func (handle *RedisHandle) Init() error {
	config := Global.Config
	for _, addr := range config.RedisAddresses {
		dialTimeout := time.Duration(config.RedisDialTimeout) * time.Second
		readTimeout := time.Duration(config.RedisReadTimeout) * time.Second
		writeTimeout := time.Duration(config.RedisWriteTimeout) * time.Second
		client := redis.NewClient(&redis.Options{
			Addr:         addr,
			DialTimeout:  dialTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			PoolSize:     config.RedisPoolSize,
		})
		_, err := client.Ping().Result()
		if err != nil {
			Global.Logger.UbLogWarning("redis init failed!")
			panic(err)
		}
		handle.Append(client)
	}
	Global.Logger.UbLogNotice("init redis OK, server num:%v", len(*handle))
	return nil
}

/**
 * @brief 获取mysql的一个实例
 * @param
 * @return
 * @author chinahbcq@qq.com
 * @date 2016-04-22 14:29:28
 */
func (handle *MysqlHandle) GetClient() (*sql.DB, bool) {
	idx := rand.Intn(len(*handle))
	return (*handle)[idx], true
}

/**
 * @brief 给mysql handle增加一个client实例
 * @param
 * @return
 * @author chinahbcq@qq.com
 * @date 2016-04-22 17:25:22
 */
func (handle *MysqlHandle) Append(client *sql.DB) {
	*handle = append(*handle, client)
}

/**
 * @brief 初始化mysql
 * @param
 * @return
 * @author chinahbcq@qq.com
 * @date 2016-04-22 14:18:22
 */
func (handle *MysqlHandle) Init() error {
	config := Global.Config
	for _, addr := range config.MysqlAddresses {
		mysqlUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", config.MysqlUsername, config.MysqlPassword, addr, config.MysqlDB)
		db, err := sql.Open("mysql", mysqlUrl)
		db.SetMaxIdleConns(0)
		(*handle).Append(db)
		if err != nil {
			Global.Logger.UbLogWarning("mysql init failed!")
			panic(err)
		}
	}
	Global.Logger.UbLogNotice("init mysql OK, server num:%v", len(*handle))
	return nil
}
