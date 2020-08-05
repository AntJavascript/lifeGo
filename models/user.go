package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 用户表字段
type AdminUser struct {
	Id       int
	UserName string
	Password string
	Email    string
	PraentId int
	Created  time.Time
	Phone    string
	Face     string
}

//  登录框用户输入字段
type PostUser struct {
	UserName string
	Password string
}

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Post() {
	o := orm.NewOrm()
	// 解析用户提交的数据
	postParams := &PostUser{}
	json.Unmarshal(c.Ctx.Input.RequestBody, postParams)
	// 构造数据查询条件
	user := AdminUser{UserName: postParams.UserName, Password: postParams.Password}

	err := o.Read(&user, "UserName", "Password")
	res := make(map[string]interface{})
	user.Password = ""
	if err == orm.ErrNoRows {
		res["data"] = ""
		res["code"] = 401
		res["msg"] = "查询不到"
	} else if err == orm.ErrMissPK {
		res["data"] = ""
		res["code"] = 401
		res["msg"] = "找不到主键"
	} else if err == nil {
		res["data"] = user
		res["code"] = 200
		res["msg"] = ""
		c.SetSession("loginuser", time.Now().UnixNano()) // 设置登录session
		fmt.Println(time.Now().UnixNano())
		fmt.Println(err)

	} else {
		res["data"] = user
		res["code"] = 500
		res["msg"] = err
	}
	c.Data["json"] = res
	c.ServeJSON()
}
