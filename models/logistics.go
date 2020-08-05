package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type LogisticsController struct {
	beego.Controller
}

// 物流表字段
type Logistics struct {
	Id            int       // id
	LogisticsName string    // 物流名称
	CreateTime    time.Time // 创建时间
}

// 获取物流列表
func (c *LogisticsController) GetList() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 获取limit ，每页需要展示的数量
	var pasize = c.GetString("limit")
	limit, err := strconv.ParseInt(pasize, 10, 64)
	// 当前处于第几页
	var page = c.GetString("page")
	currentPage, err := strconv.ParseInt(page, 10, 64)

	var list []Logistics
	num, err := o.QueryTable("Logistics").Limit(limit, currentPage-1).All(&list)

	if err != nil {
		fmt.Println(err)
		res["data"] = ""
		res["code"] = 401
		res["msg"] = err
	} else {
		res["data"] = list
		res["total"] = num
		res["code"] = 200
	}
	c.Data["json"] = res
	c.ServeJSON()
}

// 新增物流
func (c *LogisticsController) AddPost() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 构造数据查询条件
	var params Logistics
	data := c.Ctx.Input.RequestBody
	//json数据封装到user对象中
	err := json.Unmarshal(data, &params)
	fmt.Println(params)

	if err != nil {
		res["data"] = err
		res["code"] = 400
		res["msg"] = "参数解析错误"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if params.LogisticsName == "" {
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

// 修改物流
func (c *LogisticsController) EditPost() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 构造数据查询条件
	var params Logistics
	data := c.Ctx.Input.RequestBody
	//json数据封装到user对象中
	err := json.Unmarshal(data, &params)
	fmt.Println(params)

	if err != nil {
		res["data"] = err
		res["code"] = 400
		res["msg"] = "参数解析错误"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if params.LogisticsName == "" {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = "参数不存在"
	} else {
		_, err := o.Update(&params, "LogisticsName")
		if err != nil {
			fmt.Println(err)
			res["data"] = ""
			res["code"] = 401
			res["msg"] = err
		} else {
			res["data"] = ""
			res["code"] = 200
			res["msg"] = "修改成功"
		}
	}

	c.Data["json"] = res
	c.ServeJSON()
}
