package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 商品属性vaue表字段
type AttributeValue struct {
	Id             int
	AttributeId    int
	AttributeValue string
	CreateTime     time.Time
}

type AttributeValueController struct {
	beego.Controller
}

func (c *AttributeValueController) GetList() {
	res := make(map[string]interface{})
	o := orm.NewOrm()

	// 获取limit ，每页需要展示的数量
	var pasize = c.GetString("limit")
	limit, _ := strconv.ParseInt(pasize, 10, 64)
	// 当前处于第几页
	var page = c.GetString("page")
	currentPage, _ := strconv.ParseInt(page, 10, 64)

	var list []AttributeValue
	var json = make([]interface{}, 0)
	_, err := o.QueryTable("AttributeValue").Limit(limit, (currentPage-1)*limit).All(&list)
	total, totalErr := o.QueryTable("AttributeValue").Count()
	for _, value := range list {
		json = append(json, value)
	}
	if err != nil {
		res["data"] = ""
		res["code"] = 401
		res["msg"] = err
	} else {
		res["data"] = list
		res["total"] = total
		res["code"] = 200
		res["msg"] = totalErr
	}
	c.Data["json"] = res
	c.ServeJSON()

}

func (c *AttributeValueController) AddPost() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 构造数据查询条件
	var params AttributeValue
	data := c.Ctx.Input.RequestBody
	//json数据封装到user对象中
	err := json.Unmarshal(data, &params)
	if err != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = "参数解析错误"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	ParentIdIntid := strconv.Itoa(params.AttributeId)
	if params.AttributeValue == "" || ParentIdIntid == "" {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = "参数不存在"
	} else {
		params.CreateTime = time.Now()
		id, err := o.Insert(&params)
		if err != nil {
			fmt.Println(err)
			res["data"] = ""
			res["code"] = 401
			res["msg"] = err
		} else {
			res["data"] = id
			res["code"] = 200
			res["msg"] = "新增成功"
		}
	}

	c.Data["json"] = res
	c.ServeJSON()
}
