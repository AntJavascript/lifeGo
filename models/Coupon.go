package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Coupon struct {
	Id         int
	Title      string
	Desc       string
	MinAccount float64
	Price      float64
	StartTime  time.Time
	EndTime    time.Time
}
type CouponController struct {
	beego.Controller
}

// 优惠券列表
func (c *CouponController) GetList() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 获取limit ，每页需要展示的数量
	var pasize = c.GetString("limit")
	limit, err := strconv.ParseInt(pasize, 10, 64)
	// 当前处于第几页
	var page = c.GetString("page")
	currentPage, err := strconv.ParseInt(page, 10, 64)

	var list []Coupon
	var json = make([]interface{}, 0)
	num, err := o.QueryTable("Coupon").OrderBy("-id").Limit(limit, currentPage-1).All(&list)
	for _, value := range list {
		json = append(json, value)
	}
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

// 新增优惠券
func (c *CouponController) AddPost() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 构造数据查询条件
	var params Coupon
	data := c.Ctx.Input.RequestBody
	//json数据封装到params对象中
	err := json.Unmarshal(data, &params)
	if err != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
	} else {
		// 插入数据库, 返回id和err错误
		id, err := o.Insert(&params)
		if err != nil {
			res["data"] = ""
			res["code"] = 401
			res["msg"] = err
		} else {
			res["data"] = id
			res["code"] = 200
			res["msg"] = "新增成功"
		}
	}
	// 返回json格式给前端界面
	c.Data["json"] = res
	c.ServeJSON()
}
