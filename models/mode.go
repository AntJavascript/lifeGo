package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func RegisterDB() {
	//注册驱动
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//2、数据库配置
	dbHost := beego.AppConfig.String("db.host")
	dbPort := beego.AppConfig.String("db.port")
	dbDataBase := beego.AppConfig.String("db.database")
	dbUserName := beego.AppConfig.String("db.username")
	dbPwd := beego.AppConfig.String("db.pwd")
	//3、数据库连接
	conn := dbUserName + ":" + dbPwd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbDataBase + "?charset=utf8"
	//注册默认数据库
	orm.RegisterDataBase("default", "mysql", conn)

}
