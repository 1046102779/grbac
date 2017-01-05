package conf

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// db::mysql
	DBHost    string
	DBPort    int
	DBUser    string
	DBPawd    string
	DBName    string
	DBCharset string
	DBTimeLoc string
	DBMaxIdle int
	DBMaxConn int
	DBDebug   bool
)

func init() {
	var (
		err error
	)
	DBHost = strings.TrimSpace(beego.AppConfig.String("db::host"))
	if "" == DBHost {
		panic("app parameter `db::host` empty")
	}

	DBPort, err = beego.AppConfig.Int("db::port")
	if err != nil {
		panic("app parameter `db::port` error")
	}
	DBUser = strings.TrimSpace(beego.AppConfig.String("db::user"))
	if "" == DBUser {
		panic("app parameter `db::user` empty")
	}

	DBPawd = strings.TrimSpace(beego.AppConfig.String("db::pawd"))
	if "" == DBPawd {
		panic("app parameter `db::pawd` empty")
	}

	DBName = strings.TrimSpace(beego.AppConfig.String("db::name"))
	if "" == DBName {
		panic("app parameter `db::name` empty")
	}

	DBCharset = strings.TrimSpace(beego.AppConfig.String("db::charset"))
	if "" == DBCharset {
		panic("app parameter `db::charset` empty")
	}

	DBTimeLoc = strings.TrimSpace(beego.AppConfig.String("db::time_loc"))
	if "" == DBTimeLoc {
		panic("app parameter `db::time_loc` empty")
	}

	DBMaxIdle, err = beego.AppConfig.Int("db::max_idle")
	if err != nil {
		panic("app parameter `db::max_idle` error")
	}

	DBMaxConn, err = beego.AppConfig.Int("db::max_conn")
	if err != nil {
		panic("app parameter `db::max_conn` error")
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&loc=%s", DBUser, DBPawd, DBHost, DBPort, DBName, DBCharset, url.QueryEscape(DBTimeLoc))

	err = orm.RegisterDataBase("default", "mysql", dataSourceName, DBMaxIdle, DBMaxConn)
	if err != nil {
		panic("err: " + err.Error())
	}
	// orm debug
	DBDebug, err = beego.AppConfig.Bool("dev::debug")
	if err != nil {
		panic("app parameter `dev::debug` error:" + err.Error())
	}
	if DBDebug {
		orm.Debug = true
	}

	fmt.Println("conf init end")

	return
}
