package models

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 发送验证码手机号
type SendMessage struct {
	Id    int
	Phone string
	Code  int
	Type  int
}
type SendMessageController struct {
	beego.Controller
}

// 验证手机号和验证码是否匹配
func (c *SendMessageController) VerifyCode() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	var params SendMessage // 请求数据
	var sqlData SendMessage
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	sqlErr := json.Unmarshal(c.Ctx.Input.RequestBody, &sqlData)
	fmt.Println(err, sqlErr)
	if params.Phone == "" {
		res["code"] = 400
		res["msg"] = "手机号码不能为空"
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if params.Code == 0 {
		res["code"] = 400
		res["msg"] = "验证码不能为空"
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	sqlData.Phone = params.Phone
	o.QueryTable("send_message").Filter("phone", params.Phone).OrderBy("-id").Limit(1).All(&sqlData)
	fmt.Println(params)
	fmt.Println(sqlData)
	if params.Phone != sqlData.Phone {
		res["code"] = 400
		res["msg"] = "手机号码不正确"
		res["data"] = ""
	} else if params.Code != sqlData.Code {
		res["code"] = 400
		res["msg"] = "验证码不正确"
		res["data"] = ""
	} else {
		res["code"] = 200
		res["msg"] = "ok"
		res["data"] = ""
	}
	c.Data["json"] = res
	c.ServeJSON()

}
