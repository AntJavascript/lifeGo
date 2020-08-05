package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 商品属性key表字段
type AttributeKey struct {
	Id            int
	CategoryId    int
	AttributeName string
	CreateTime    time.Time
}

// 返回前端结构
type CallBackJson struct {
	AttributeKey
	ValueList interface{}
}
type AttributeKeyController struct {
	beego.Controller
}

// 查找key对应的value
func findValue(attributeKey AttributeKey, keyId int) CallBackJson {
	fmt.Println(keyId)
	o := orm.NewOrm()
	var AttributeValueJson []AttributeValue
	var callBackJson CallBackJson

	o.QueryTable("AttributeValue").Filter("AttributeId", keyId).All(&AttributeValueJson) // 获取一级商品分类
	fmt.Println(AttributeValueJson)
	callBackJson.AttributeKey = attributeKey
	callBackJson.ValueList = AttributeValueJson

	return callBackJson
}

func (c *AttributeKeyController) GetList() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	var list []AttributeKey
	var json = make([]interface{}, 0)
	num, err := o.QueryTable("AttributeKey").All(&list)
	for _, value := range list {
		json = append(json, findValue(value, value.Id))
	}
	if err != nil {
		res["data"] = ""
		res["code"] = 401
		res["msg"] = err
	} else {
		res["data"] = &json
		res["total"] = num
		res["code"] = 200
	}
	c.Data["json"] = res
	c.ServeJSON()

}

func (c *AttributeKeyController) AddPost() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 构造数据查询条件
	var params AttributeKey
	data := c.Ctx.Input.RequestBody
	//json数据封装到user对象中
	err := json.Unmarshal(data, &params)
	fmt.Println(err)
	if params.AttributeName == "" {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = "AttributeName 不存在"
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
