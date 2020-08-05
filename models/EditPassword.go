package models

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 请求的参数结构
type EditpasswordParams struct {
	Phone           string
	Password        string
	ConfirmPassword string
}
type EditpasswordController struct {
	beego.Controller
}

func (c *EditpasswordController) EditPassword() {
	o := orm.NewOrm()
	res := make(map[string]interface{})
	// 数据库字段结构
	var params User
	paramsErr := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	// 前端提交的数据结构
	var EditpasswordParams EditpasswordParams
	RegisterParamsErr := json.Unmarshal(c.Ctx.Input.RequestBody, &EditpasswordParams)

	if paramsErr == nil && RegisterParamsErr == nil {
		if EditpasswordParams.Phone == "" {
			res["code"] = 400
			res["msg"] = "手机号码不能为空"
			res["data"] = ""
			c.Data["json"] = res
			c.ServeJSON()
			return
		}
		if EditpasswordParams.Password == "" {
			res["code"] = 400
			res["msg"] = "密码不能为空"
			res["data"] = ""
			c.Data["json"] = res
			c.ServeJSON()
			return
		}
		if EditpasswordParams.ConfirmPassword == "" {
			res["code"] = 400
			res["msg"] = "确认密码不能为空"
			res["data"] = ""
			c.Data["json"] = res
			c.ServeJSON()
			return
		}
		if EditpasswordParams.ConfirmPassword != EditpasswordParams.Password {
			res["code"] = 400
			res["msg"] = "密码不一致"
			res["data"] = ""
			c.Data["json"] = res
			c.ServeJSON()
			return
		}
		var sqlData User
		// 用于查询是否已注册
		sqlDataErr := json.Unmarshal(c.Ctx.Input.RequestBody, &sqlData)
		fmt.Println(sqlDataErr)
		o.QueryTable("user").Filter("phone", params.Phone).Limit(1).All(&sqlData)

		h := md5.New()
		password := []byte(params.Password)
		// md5加密
		md5str := fmt.Sprintf("%x", h.Sum(password))
		params.Password = md5str
		params.Id = sqlData.Id
		fmt.Println(params)
		// 获取当前时间
		params.CreteTime = time.Now()
		id, err := o.Update(&params, "password")
		if err != nil {
			res["code"] = 401
			res["msg"] = "密码修改失败"
			res["data"] = ""
			c.Data["json"] = res
			c.ServeJSON()
		} else {
			res["code"] = 200
			res["msg"] = "密码修改成功"
			res["data"] = id
			c.Data["json"] = res
			c.ServeJSON()
		}

	} else {
		res["code"] = 400
		res["msg"] = "参数解析错误"
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
	}
}
