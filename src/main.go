/***************************************************************************
 *
 * Copyright (c)2016 SmartApi. All Rights Reserved
 *
 **************************************************************************/


/**
 * @file main.go
 * @author chenqian (chinahbcq@qq.com)
 * @date 2016-04-29 15:27:13
 * @brief 
 *
 **/

package main

import (
	//"io"
	"io/ioutil"
	"net/http"
	//"path"
	"action"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ublog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)
import (
	"utils"
)

func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
			}
		}()
		fn(w, r)
	}
}
func loadSysConfigs(file string) error {
	configStr, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configStr, &utils.Global.Config)
	if err != nil {
		panic(err)
	}
	return nil
}
func loadErrorConfigs(file string) error {
	configStr, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	//var r interface{}
	//err = json.Unmarshal(configStr, &r)
	//codes, ok := r.(map[string] interface{})
	//if ok {
	//    for j,k := range codes {
	//        log.Println(j, k)
	//    }
	//}

	err = json.Unmarshal(configStr, &utils.ErrorMap)
	if err != nil {
		return err
	}
	return nil
}

func initDB() error {
	//初始化mysql
	err := utils.Global.MysqlHandles.Init()
	if err != nil {
		utils.Global.Logger.UbLogWarning("init mysql failed:%s", err.Error())
		return err
	}
	utils.Global.Logger.UbLogNotice("init mysql OK")
	//初始化redis
	err = utils.Global.RedisHandles.Init()
	if err != nil {
		utils.Global.Logger.UbLogWarning("init redis failed:%s", err.Error())
		return err
	}
	utils.Global.Logger.UbLogNotice("init redis OK")
	return nil
}

func initLogger() {
	config := utils.Global.Config
	ubLogInfo := &ublog.UbLogInfo{
		ChannelLen:     config.LogChannelLen,
		FlushThreadNum: config.LogFlushThreadNum,
		FlushInterval:  config.LogFlushInterval,
	}
	if "" != config.LogFile {
		ublog.UbLogFdPool[3] = &ublog.UbLog{}
		ublog.UbLogFdPool[3].Init(config.LogDir,
			config.LogFile,
			config.LogLevel,
			ubLogInfo,
			false)
		utils.Global.Logger = ublog.UbLogFdPool[3]
	} else {
		utils.Global.Logger = ublog.UbLogFdPool[2]
	}
	go func() {
		for {
			utils.Global.Logger.CheckLogFile()
			time.Sleep(1000 * time.Millisecond)
		}
	}()
}

func main() {
	//1 读取系统配置
	loadSysConfigs("conf/smartapi.conf")

	//2 初始化log
	initLogger()

	configs := utils.Global.Config
	logger := utils.Global.Logger
	logger.UbLogNotice("init ublog OK")

	//3 读取错误码配置
	err := loadErrorConfigs("conf/error_code.conf")
	if err != nil {
		panic(err)
	}
	logger.UbLogNotice("load error_code.conf OK")

	//4 初始化数据库
	err = initDB()
	if err != nil {
		panic(err)
	}
	logger.UbLogNotice("init all dbs OK")
	runtime.GOMAXPROCS(configs.ProcessNum)

	//5 开启服务
	go func() {
		mux := mux.NewRouter()
		mux.HandleFunc("/api/1.0/users.{name:[a-zA-Z]+}", safeHandler(action.Users))

		err := http.ListenAndServe(configs.ListenAddress, mux)
		if err != nil {
			logger.UbLogWarning("start service fail %s", err.Error())
			panic(err)
		}
	}()
	logger.UbLogNotice("start service at %s OK", configs.ListenAddress)

	//6 响应退出信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	logger.UbLogNotice("service exist!")
}
