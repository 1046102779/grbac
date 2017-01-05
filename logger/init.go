package logger

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var (
	LogId string
)

var (
	Logger     *logs.BeeLogger
	LoggerSMTP *logs.BeeLogger
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

	loggerFileFilename = strings.TrimSpace(beego.AppConfig.String("logger_file::filename"))
	if "" == loggerFileFilename {
		loggerFileFilename = "log/" + beego.BConfig.AppName + ".log"
	}

	loggerFileMaxlinesStr := strings.TrimSpace(beego.AppConfig.String("logger_file::maxlines"))
	if "" == loggerFileMaxlinesStr {
		loggerFileMaxlines = 1000000
	} else {
		loggerFileMaxlines, err = beego.AppConfig.Int("logger_file::maxlines")
		if err != nil {
			panic("app conf `logger_file::maxlines` error:" + err.Error())
		}
	}

	loggerFileMaxsizeStr := strings.TrimSpace(beego.AppConfig.String("logger_file::maxsize"))
	if "" == loggerFileMaxsizeStr {
		loggerFileMaxsize = 256 * 1024 * 1024
	} else {
		loggerFileMaxsize, err = beego.AppConfig.Int64("logger_file::maxsize")
		if err != nil {
			panic("app conf `logger_file::maxsize` error:" + err.Error())
		}
		loggerFileMaxsize = loggerFileMaxsize * 1024 * 1024
	}

	loggerFileDailyStr := strings.TrimSpace(beego.AppConfig.String("logger_file::daily"))
	if "" == loggerFileDailyStr {
		loggerFileDaily = true
	} else {
		loggerFileDaily, err = beego.AppConfig.Bool("logger_file::daily")
		if err != nil {
			panic("app conf `logger_file::daily` error:" + err.Error())
		}
	}

	loggerFileMaxdaysStr := strings.TrimSpace(beego.AppConfig.String("logger_file::maxdays"))
	if "" == loggerFileMaxdaysStr {
		loggerFileMaxdays = 7
	} else {
		loggerFileMaxdays, err = beego.AppConfig.Int("logger_file::maxdays")
		if err != nil {
			panic("app conf `logger_file::maxdays` error:" + err.Error())
		}
	}

	loggerFileRotateStr := strings.TrimSpace(beego.AppConfig.String("logger_file::rotate"))
	if "" == loggerFileRotateStr {
		loggerFileRotate = true
	} else {
		loggerFileRotate, err = beego.AppConfig.Bool("logger_file::rotate")
		if err != nil {
			panic("app conf `logger_file::rotate` error:" + err.Error())
		}
	}

	loggerFileLevelStr := strings.TrimSpace(beego.AppConfig.String("logger_file::level"))
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
			panic("app conf `logger_file::level error: not defined value")
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
	logFuncCallEnable, err := beego.AppConfig.Bool("logger_file::log_func_call_enable")
	if err != nil {
		logFuncCallEnable = true
	}
	Logger = logs.NewLogger(10000)
	Logger.EnableFuncCallDepth(logFuncCallEnable)
	Logger.SetLogFuncCallDepth(2)

	Logger.SetLogger("file", loggerConf)
}
