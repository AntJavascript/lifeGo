package models

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// about表字段
type About struct {
	Id      int
	Content string
}

type AboutController struct {
	beego.Controller
}

func (c *AboutController) GetAuout() {
	o := orm.NewOrm()
	res := make(map[string]interface{})

	var list []About
	_, err := o.QueryTable("About").Limit(1).All(&list)
	if err != nil {
		res["data"] = ""
		res["msg"] = "查询错误"
		res["code"] = 401
	} else {
		res["data"] = list
		res["code"] = 200
	}

	c.Data["json"] = res
	c.ServeJSON()
}

func (c *AboutController) EditAbout() {
	o := orm.NewOrm()
	res := make(map[string]interface{})

	var params About
	// 用于查询是否已有记录
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		res["data"] = ""
		res["msg"] = "参数解析错误"
		res["code"] = 400
	} else {
		if params.Id != 0 {
			id, err := o.Update(&params, "content")
			fmt.Println(id)
			if err != nil {
				res["data"] = id
				res["msg"] = "修改失败"
				res["code"] = 401
			} else {
				res["data"] = id
				res["msg"] = "修改成功"
				res["code"] = 200
			}
		} else {
			id, err := o.Insert(&params)
			if err != nil {
				res["data"] = ""
				res["msg"] = "新增失败"
				res["code"] = 401
			} else {
				res["data"] = id
				res["msg"] = "新增成功"
				res["code"] = 200
			}
		}
	}

	c.Data["json"] = res
	c.ServeJSON()
}
