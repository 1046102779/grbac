package conf

import (
	"fmt"
	"net/url"
	"strings"

	redis "gopkg.in/redis.v5"

	"github.com/1046102779/grbac/common/utils"
	"github.com/smallnest/rpcx"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	AccountClient     *rpcx.Client
	EtcdAddr, RpcAddr string
	Servers           []string
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

	// redis
	Redis__Address  string
	Redis__Password string
	Redis__Db       int

	Redis__Client *utils.RedisV5Client
)

func initRpcEnv() {
	EtcdAddr = strings.TrimSpace(beego.AppConfig.String("etcd::address"))
	RpcAddr = strings.TrimSpace(beego.AppConfig.String("rpc::address"))
	if EtcdAddr == "" || RpcAddr == "" {
		panic("params `etcd::address || rpc::address` empty")
	}
	serverTemp := beego.AppConfig.String("rpc::servers")
	Servers = strings.Split(serverTemp, ",")
	return
}
func initRedis() {
	var (
		err error
	)

	Redis__Address = beego.AppConfig.String("redis::address")
	if "" == Redis__Address {
		panic("parameter `redis::address` empty")
	}

	Redis__Password = beego.AppConfig.String("redis::password")
	if "" == Redis__Password {
		panic("parameter `redis::password` empty")
	}

	Redis__Db, err = beego.AppConfig.Int("redis::db")
	if err != nil {
		panic("parameter `redis::db` error")
	}

	Redis__Client = &utils.RedisV5Client{
		Options: &redis.Options{
			Addr:     Redis__Address,
			Password: Redis__Password,
			DB:       Redis__Db,
		},
	}

	if err = Redis__Client.Conn(); err != nil {
		panic(err)
	}

}

func initDB() {
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

	return
}

func init() {
	initRedis()
	initDB()
	initRpcEnv()
	// orm debug
	DBDebug, err := beego.AppConfig.Bool("dev::debug")
	if err != nil {
		panic("app parameter `dev::debug` error:" + err.Error())
	}
	if DBDebug {
		orm.Debug = true
	}

	fmt.Println("conf init end")
}
