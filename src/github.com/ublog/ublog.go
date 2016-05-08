/**
 * @file ublog.go
 * @author	qiuxueda@baidu.com
 * @brief   仿照php-ublog库制作的golang日志库
 *
 * @history
 *	v0.1	输出日志到标准输出及错误输出，搭配comlog_apache可实现日志的高级打印及切分等
 *
 **/
package ublog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"time"
	//	"github.com/glog"
)

var UbLogDefaultChannelLen = 100000
var UbLogDefaultLogLevel = uint32(16)

var UbLogFdPool map[uint32]*UbLog

type UbLog struct {
	logPath       string
	logName       string
	logLevel      uint32
	isDio         bool
	LogChannel    chan string
	WfLogChannel  chan string
	fd            io.WriteCloser
	wfFd          io.WriteCloser
	channelLen    uint32
	flushInterval uint32
}

func init() {
	UbLogFdPool = make(map[uint32]*UbLog)
	UbLogFdPool[1] = &UbLog{}
	UbLogFdPool[1].Init("", "", UbLogDefaultLogLevel, &UbLogInfo{}, false)
	UbLogFdPool[2] = &UbLog{}
	UbLogFdPool[2].Init("", "", UbLogDefaultLogLevel, &UbLogInfo{LogFd: os.Stderr}, false)
}

type UbLogInfo struct {
	ChannelLen     uint32
	LogFd          io.WriteCloser //如果不为空，则优先使用该fd而不是打开文件
	FlushThreadNum uint32
	FlushInterval  uint32
}

func (ul *UbLog) Init(logPath string, logName string, logLevel uint32, logInfo *UbLogInfo, isDio bool) (err bool) {
	ul.logLevel = logLevel
	if logLevel <= 0 {
		logLevel = UbLogDefaultLogLevel
	}
	if logInfo.ChannelLen <= 0 {
		logInfo.ChannelLen = 10000
	}
	if nil != logInfo.LogFd {
		ul.fd = logInfo.LogFd
	} else if "" == logName {
		ul.fd = os.Stdout
	} else {
		os.MkdirAll(logPath, 0755)
		ul.fd, _ = os.OpenFile(logPath+"/"+logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		ul.wfFd, _ = os.OpenFile(logPath+"/"+logName+".wf", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		ul.logPath = logPath
		ul.logName = logName
	}
	ul.channelLen = logInfo.ChannelLen
	ul.LogChannel = make(chan string, ul.channelLen)
	ul.WfLogChannel = make(chan string, ul.channelLen)
	if logInfo.FlushThreadNum <= 0 {
		logInfo.FlushThreadNum = 1
	}
	if logInfo.FlushInterval <= 0 {
		logInfo.FlushInterval = 1
	}
	ul.flushInterval = logInfo.FlushInterval
	for num := uint32(0); num < logInfo.FlushThreadNum; num++ {
		go ul.FlushLog()
	}
	return true
}

func checkFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//only for removed files
func (ul *UbLog) CheckLogFile() {
	ul.updateLogFile("")
	ul.updateLogFile(".wf")
}

func (ul *UbLog) updateLogFile(f string) {
	var fd io.WriteCloser
	var ofd io.WriteCloser
	var err error

	fname := ul.logPath + "/" + ul.logName + f
	if checkFileExist(fname) {
		return
	}

	//fmt.Println("log not exist: ", fname)
	fd, err = os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if nil != err {
		//fmt.Println("open new log err: ", err)
		return
	}

	//fmt.Println("open new log: ", fd)

	switch f {
	case ".wf":
		ofd = ul.wfFd
		ul.wfFd = fd
	default:
		ofd = ul.fd
		ul.fd = fd
	}

	if nil != ofd {
		//fmt.Println("close old log: ", ofd)
		ofd.Close()
	}
}

func (ul *UbLog) FlushLog() {
	var logBuffer bytes.Buffer
	var wfLogBuffer bytes.Buffer
	for {
		for i := uint32(0); i < ul.flushInterval; i++ {
			select {
			case logMessage := <-ul.LogChannel:
				logBuffer.Write([]byte(logMessage))
			case logMessage := <-ul.WfLogChannel:
				wfLogBuffer.Write([]byte(logMessage))
			}

		}
		ul.fd.Write(logBuffer.Bytes())
		logBuffer.Reset()

		ul.wfFd.Write(wfLogBuffer.Bytes())
		wfLogBuffer.Reset()
	}
}

func (ul *UbLog) UbLogWriteLog(logLevel uint32, format string, a ...interface{}) (n int, err error) {
	if logLevel <= ul.logLevel {
		nowTime := time.Now()
		logTime := nowTime.Format("2006-01-02 15:04:05.000000")
		funcName, _, lineNum, _ := runtime.Caller(2)
		logChannel := ul.LogChannel

		var logLevelName = ""
		switch logLevel {
		case 8:
			logLevelName = "TRACE"
			//glog.Infoln(logLevelName, fmt.Sprintf(format, a...))
			break
		case 16:
			logLevelName = "DEBUG"
			//glog.Infoln(logLevelName, fmt.Sprintf(format, a...))
			break
		case 4:
			logLevelName = "NOTICE"
			//glog.Infoln(logLevelName, fmt.Sprintf(format, a...))
			break
		case 2:
			logLevelName = "WARNING"
			logChannel = ul.WfLogChannel
			//glog.Warningln(logLevelName, fmt.Sprintf(format, a...))
			break
		case 1:
			logLevelName = "FATAL"
			logChannel = ul.WfLogChannel
			//glog.Errorln(logLevelName, fmt.Sprintf(format, a...))
			break
		}

		logChannel <- fmt.Sprintln(logLevelName + " " + logTime + " " + runtime.FuncForPC(funcName).Name() + ":" + strconv.Itoa(lineNum) + " " + fmt.Sprintf(format, a...))
	}
	return 0, nil
}

func (ul *UbLog) UbLogTrace(format string, a ...interface{}) (n int, err error) {
	return ul.UbLogWriteLog(16, format, a...)
}

func (ul *UbLog) UbLogDebug(format string, a ...interface{}) (n int, err error) {
	return ul.UbLogWriteLog(8, format, a...)
}

func (ul *UbLog) UbLogNotice(format string, a ...interface{}) (n int, err error) {
	return ul.UbLogWriteLog(4, format, a...)
}

func (ul *UbLog) UbLogWarning(format string, a ...interface{}) (n int, err error) {
	return ul.UbLogWriteLog(2, format, a...)
}

func (ul *UbLog) UbLogFatal(format string, a ...interface{}) (n int, err error) {
	return ul.UbLogWriteLog(1, format, a...)
}
