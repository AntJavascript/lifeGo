package models

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 登录记录表字段
type Login struct {
	Id        int
	Phone     string
	LoginTime time.Time
	LastTime  time.Time
}

// 前端参数字段
type LoginParams struct {
	Phone    string
	Password string
}

type AppLoginController struct {
	beego.Controller
}

func (c *AppLoginController) Login() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	//	var loginParams Login
	var AjaxParams LoginParams
	// 解析前端参数
	paramsErr := json.Unmarshal(c.Ctx.Input.RequestBody, &AjaxParams)
	fmt.Println(AjaxParams)
	var sqlData User
	o.QueryTable("user").Filter("phone", AjaxParams.Phone).Limit(1).All(&sqlData)
	fmt.Println(sqlData)
	h := md5.New()
	password := []byte(AjaxParams.Password)
	// md5加密
	md5str := fmt.Sprintf("%x", h.Sum(password))
	if AjaxParams.Phone != sqlData.Phone {
		res["code"] = 400
		res["msg"] = "用户名称不存在"
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if md5str != sqlData.Password {
		res["code"] = 400
		res["msg"] = "密码错误"
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if paramsErr != nil {
		res["code"] = 400
		res["msg"] = "参数解析错误"
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if AjaxParams.Phone == "" {
		res["code"] = 400
		res["msg"] = "手机号码不能为空"
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	c.SetSession("LoginSession", sqlData)

	sessionId := c.CruSession

	res["code"] = 200
	res["msg"] = "登录成功"
	res["data"] = sessionId
	c.Data["json"] = res
	c.ServeJSON()
}

// 退出
func (c *AppLoginController) LoginOut() {
	res := make(map[string]interface{})
	c.DelSession("LoginSession")

	res["code"] = 200
	res["msg"] = ""
	res["data"] = ""
	c.Data["json"] = res
	c.ServeJSON()

}
