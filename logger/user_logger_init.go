package logger

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var (
	UserLogger *logs.BeeLogger
)

func init() {
	var (
		loggerFileFilename string
		loggerFileMaxlines int
		loggerFileMaxsize  int64
		loggerFileDaily    bool
		loggerFileMaxdays  int
		loggerFileRotate   bool
		loggerFileLevel    int
		err                error
	)

	loggerFilePath := strings.TrimSpace(beego.AppConfig.String("user_log::user_log_path"))
	if "" == loggerFilePath {
		panic("app conf `user_log::user_log_path` error:" + err.Error())
	}

	loggerFileFilename = strings.TrimSpace(beego.AppConfig.String("user_log::filename"))
	if "" == loggerFileFilename {
		loggerFileFilename = loggerFilePath + "/" + beego.BConfig.AppName + ".log"
	}

	loggerFileMaxlinesStr := strings.TrimSpace(beego.AppConfig.String("user_log::maxlines"))
	if "" == loggerFileMaxlinesStr {
		loggerFileMaxlines = 1000000
	} else {
		loggerFileMaxlines, err = beego.AppConfig.Int("user_log::maxlines")
		if err != nil {
			panic("app conf `user_log::maxlines` error:" + err.Error())
		}
	}

	loggerFileMaxsizeStr := strings.TrimSpace(beego.AppConfig.String("user_log::maxsize"))
	if "" == loggerFileMaxsizeStr {
		loggerFileMaxsize = 256 * 1024 * 1024
	} else {
		loggerFileMaxsize, err = beego.AppConfig.Int64("user_log::maxsize")
		if err != nil {
			panic("app conf `user_log::maxsize` error:" + err.Error())
		}
		loggerFileMaxsize = loggerFileMaxsize * 1024 * 1024
	}

	loggerFileDailyStr := strings.TrimSpace(beego.AppConfig.String("user_log::daily"))
	if "" == loggerFileDailyStr {
		loggerFileDaily = true
	} else {
		loggerFileDaily, err = beego.AppConfig.Bool("user_log::daily")
		if err != nil {
			panic("app conf `user_log::daily` error:" + err.Error())
		}
	}

	loggerFileMaxdaysStr := strings.TrimSpace(beego.AppConfig.String("user_log::maxdays"))
	if "" == loggerFileMaxdaysStr {
		loggerFileMaxdays = 7
	} else {
		loggerFileMaxdays, err = beego.AppConfig.Int("user_log::maxdays")
		if err != nil {
			panic("app conf `user_log::maxdays` error:" + err.Error())
		}
	}

	loggerFileRotateStr := strings.TrimSpace(beego.AppConfig.String("user_log::rotate"))
	if "" == loggerFileRotateStr {
		loggerFileRotate = true
	} else {
		loggerFileRotate, err = beego.AppConfig.Bool("user_log::rotate")
		if err != nil {
			panic("app conf `user_log::rotate` error:" + err.Error())
		}
	}

	loggerFileLevelStr := strings.TrimSpace(beego.AppConfig.String("user_log::level"))
	if "" == loggerFileLevelStr {
		loggerFileLevel = logs.LevelDebug
	} else {
		switch loggerFileLevelStr {
		case "Emergency":
			loggerFileLevel = logs.LevelEmergency
		case "Alert":
			loggerFileLevel = logs.LevelAlert
		case "Critical":
			loggerFileLevel = logs.LevelCritical
		case "Error":
			loggerFileLevel = logs.LevelError
		case "Warning":
			loggerFileLevel = logs.LevelWarning
		case "Notice":
			loggerFileLevel = logs.LevelNotice
		case "Informational":
			loggerFileLevel = logs.LevelInformational
		case "Debug":
			loggerFileLevel = logs.LevelDebug
		default:
			panic("app conf `user_log::level error: not defined value")
		}
	}

	loggerConf := fmt.Sprintf(`{
		        "filename": "%s",
		        "maxlines": %d,
		        "maxsize": %d,
		        "daily": %t,
		        "maxdays": %d,
		        "rotate": %t,
		        "level": %d
		    }`, loggerFileFilename, loggerFileMaxlines, loggerFileMaxsize, loggerFileDaily, loggerFileMaxdays, loggerFileRotate, loggerFileLevel)
	logFuncCallEnable, err := beego.AppConfig.Bool("user_log::log_func_call_enable")
	if err != nil {
		logFuncCallEnable = true
	}
	UserLogger = logs.NewLogger(10000)
	UserLogger.EnableFuncCallDepth(logFuncCallEnable)
	UserLogger.SetLogFuncCallDepth(2)

	UserLogger.SetLogger("file", loggerConf)
}
