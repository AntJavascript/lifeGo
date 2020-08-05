package models

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserCoupon struct {
	Id       int
	UserId   int
	CouponId int
	IsUse    int
}
type UserCouponController struct {
	beego.Controller
}

// 个人优惠券列表
func (c *UserCouponController) GetList() {
	res := make(map[string]interface{})
	o := orm.NewOrm()

	loginSession := c.GetSession("LoginSession")
	if loginSession == nil {
		res["code"] = 401
		res["msg"] = "请先登录"
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	var flagUser int
	if v, ok := loginSession.(User); ok {
		flagUser = v.Id
	}
	var list []UserCoupon
	num, err := o.QueryTable("UserCoupon").Filter("UserId", flagUser).All(&list)

	var user_coupon_list = make([]interface{}, 0) // 个人优惠券列表

	for _, item := range list {
		var result_data = make(map[string]interface{}, 0) // 优惠券数据
		result_data["IsUse"] = item.IsUse

		var coupon_data Coupon
		coupon_data.Id = item.CouponId
		coupon_data_err := o.Read(&coupon_data)

		if coupon_data_err == nil {
			result_data["Desc"] = coupon_data.Desc
			result_data["EndTime"] = coupon_data.EndTime
			result_data["Id"] = coupon_data.Id
			result_data["MinAccount"] = coupon_data.MinAccount
			result_data["Price"] = coupon_data.Price
			result_data["StartTime"] = coupon_data.StartTime
			result_data["Title"] = coupon_data.Title
		}
		user_coupon_list = append(user_coupon_list, result_data)

	}

	if err != nil {
		res["data"] = ""
		res["code"] = 401
		res["msg"] = err
	} else {
		res["data"] = list
		res["user_list"] = user_coupon_list
		res["total"] = num
		res["code"] = 200
	}
	c.Data["json"] = res
	c.ServeJSON()
}

// 领取优惠券
func (c *UserCouponController) ReceiveCoupon() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	loginSession := c.GetSession("LoginSession")
	if loginSession == nil {
		res["code"] = 401
		res["msg"] = "请先登录"
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	var flagUser int
	if v, ok := loginSession.(User); ok {
		flagUser = v.Id
	}
	// 构造数据查询条件
	var params UserCoupon
	//json数据封装到params对象中
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	params.UserId = flagUser
	params.IsUse = 0
	if err != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
	} else {
		var list []UserCoupon
		o.QueryTable("UserCoupon").Filter("UserId", flagUser).All(&list)
		for _, item := range list {
			if item.CouponId == params.CouponId {
				res["data"] = ""
				res["code"] = 401
				res["msg"] = "已经领取，请勿重复领取"
				// 返回json格式给前端界面
				c.Data["json"] = res
				c.ServeJSON()
				return
			}
		}
		// 插入数据库, 返回id和err错误
		id, err := o.Insert(&params)
		if err != nil {
			res["data"] = ""
			res["code"] = 401
			res["msg"] = err
		} else {
			res["data"] = id
			res["code"] = 200
			res["msg"] = "领取成功"
		}
	}
	// 返回json格式给前端界面
	c.Data["json"] = res
	c.ServeJSON()
}
