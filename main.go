package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/casbin/beego-orm-adapter"
	"github.com/casbin/casbin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/qiangjiyi/casbin_auth/filter"
	_ "github.com/qiangjiyi/casbin_auth/routers"
)

var e *casbin.Enforcer

func init() {
	// set beego logs engine
	logs.SetLogger(logs.AdapterConsole, `{"level":7,"color":true}`)
	logs.SetLogger(logs.AdapterFile, `{"filename":"./logs/casbin_auth.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	logs.SetLogFuncCallDepth(3)

	driverName := beego.AppConfig.String("db_driverName")
	dataSource := beego.AppConfig.String("db_dataSource")

	// construct a adapter base on beego-orm as casbin's policy storage
	adapter := beegoormadapter.NewAdapter(driverName, dataSource, true)
	e = casbin.NewEnforcer("./conf/rbac_model.conf", adapter, true)
}

func main() {
	// open swagger api function
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	// add permission authentication filter
	beego.InsertFilter("*", beego.BeforeRouter, filter.NewAuthorizer(e, beego.BConfig.WebConfig.Session.SessionOn))

	beego.Run()
}
