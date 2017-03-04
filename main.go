package main

import (
	"fmt"
	"time"

	"github.com/1046102779/grbac/conf"
	. "github.com/1046102779/grbac/logger"
	_ "github.com/1046102779/grbac/routers"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx"
	"github.com/smallnest/rpcx/codec"
	"github.com/smallnest/rpcx/plugin"

	"github.com/1046102779/grbac/models"
	"github.com/astaxie/beego"
)

func startRPCService(rpcAddr string, etcdAddr string, grbacServer *models.GrbacServer) {
	server := rpcx.NewServer()
	rplugin := &plugin.EtcdRegisterPlugin{
		ServiceAddress: "tcp@" + rpcAddr,
		EtcdServers:    []string{etcdAddr},
		BasePath:       fmt.Sprintf("/%s/%s", beego.BConfig.RunMode, "rpcx"),
		Metrics:        metrics.NewRegistry(),
		Services:       make([]string, 0),
		UpdateInterval: time.Minute,
	}
	rplugin.Start()
	server.PluginContainer.Add(rplugin)
	server.PluginContainer.Add(plugin.NewMetricsPlugin())
	server.RegisterName("grbacs", grbacServer, "weight=1&m=devops")
	server.ServerCodecFunc = codec.NewProtobufServerCodec
	server.Serve("tcp", rpcAddr)
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	// 1. 加载路由树数据
	if err := LoadMapTree(); err != nil {
		Logger.Error(err.Error())
	}
	// 打印初始化的树
	models.PrintTree()

	// 对于存储在Redis集群中的数据，一般情况下，只需要加载一次
	// 以后重启权限服务，不需要再次加载
	// 所以这是一个python脚本, 脚本路径：/data/home/chendonghai/data/python/grbac_init_data_mysql_to_redis.py
	// 2. 加载角色与功能的Redis存储  数据结构：SET集合
	// 3. 加载用户与角色的Redis存储  数据结构：SET集合
	// 4. 加载<用户-功能ID，实体ID>的Redis存储	  数据结构：SET集合
	//	备注：目前只有公司ID
	// 5. 加载白名单的Redis存储 数据结构：LIST列表
	// 6. 加载<域，动作，实体>列表在redis缓存中
	// 数据结构集合,  格式：<用户ID-功能ID, 实体SET>
	/*
		if err := LoadEntity(); err != nil {
		}
	*/
	if _, err := models.LoadWhiteList(0, 300); err != nil {
		Logger.Error(err.Error())
	}
	fmt.Println("main starting...")
	go startRPCService(conf.RpcAddr, conf.EtcdAddr, &models.GrbacServer{})

	beego.Run()
}
