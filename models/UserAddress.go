package models

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 收货地址表字段
type UserAddress struct {
	Id        int
	UserId    int
	IsDefault int
	Province  string
	City      string
	Area      string
	Detail    string
	Phone     string
	Consignee string
}
type PostParams struct {
	Id     int
	UserId int
}
type UserAddressController struct {
	beego.Controller
}

// 获取个人收货地址
func (c *UserAddressController) GetUserAddressList() {
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
	var list []UserAddress
	_, err := o.QueryTable("UserAddress").Filter("user_id", flagUser).All(&list)
	if err != nil {
		res["code"] = 401
		res["code"] = 200
		res["msg"] = err
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res["code"] = 200
		res["msg"] = err
		res["data"] = list
		c.Data["json"] = res
		c.ServeJSON()
	}
}

// 添加个人收货地址
func (c *UserAddressController) AddUserAddress() {
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
	var params UserAddress
	// 解析前端参数
	paramsErr := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if paramsErr != nil {
		res["code"] = 401
		res["msg"] = paramsErr
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	params.UserId = flagUser
	if params.IsDefault == 1 {
		var list []UserAddress
		_, rowErr := o.QueryTable("UserAddress").Filter("UserId", flagUser).All(&list)
		// 把所有的默认收货地址变成0
		for _, item := range list {
			item.IsDefault = 0
			_, err := o.Update(&item, "IsDefault")
			fmt.Println(err)
		}

		fmt.Println(list, rowErr)

	}

	id, err := o.Insert(&params)
	if err != nil {
		res["code"] = 401
		res["msg"] = err
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res["code"] = 200
		res["msg"] = err
		res["data"] = id
		c.Data["json"] = res
		c.ServeJSON()
	}
}

// 删除个人收货地址
func (c *UserAddressController) DeleteUserAddress() {
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

	var params PostParams
	// 解析前端参数
	paramsErr := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if paramsErr != nil {
		res["code"] = 401
		res["msg"] = paramsErr
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	params.UserId = flagUser
	fmt.Println(params)
	num, err := o.QueryTable("UserAddress").Filter("UserId", flagUser).Filter("Id", params.Id).Delete()
	if err != nil {
		res["code"] = 401
		res["msg"] = err
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res["code"] = 200
		res["msg"] = "删除成功"
		res["data"] = num
		c.Data["json"] = res
		c.ServeJSON()
	}
}

// 个人收货地址详情
func (c *UserAddressController) UserAddressDetail() {
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

	var params PostParams

	// 获取ajax提交的参数
	paramsId := c.Ctx.Input.Param(":id")
	id, paramsErr := strconv.Atoi(paramsId)
	fmt.Println(id)
	if paramsErr != nil {
		res["code"] = 401
		res["msg"] = paramsErr
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	params.Id = id
	params.UserId = flagUser
	fmt.Println(params)
	var resultData UserAddress
	_, err := o.QueryTable("UserAddress").Filter("UserId", flagUser).Filter("Id", params.Id).All(&resultData)
	if err != nil {
		res["code"] = 401
		res["msg"] = err
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res["code"] = 200
		res["msg"] = ""
		res["data"] = resultData
		c.Data["json"] = res
		c.ServeJSON()
	}
}

// 修改个人收货地址
func (c *UserAddressController) EditUserAddress() {
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
	var params UserAddress
	// 解析前端参数
	paramsErr := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	fmt.Println(params)
	if paramsErr != nil {
		res["code"] = 401
		res["msg"] = paramsErr
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	params.UserId = flagUser
	fmt.Println("+++++++++++++++++++++++++")
	fmt.Println(params)
	if params.IsDefault == 1 {
		var list []UserAddress
		_, rowErr := o.QueryTable("UserAddress").Filter("UserId", flagUser).All(&list)
		fmt.Println(rowErr)
		// 把所有的默认收货地址变成0
		for _, item := range list {
			item.IsDefault = 0
			_, err := o.Update(&item, "IsDefault")
			fmt.Println(err)
		}
	}

	id, err := o.Update(&params)
	if err != nil {
		res["code"] = 401
		res["msg"] = err
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res["code"] = 200
		res["msg"] = err
		res["data"] = id
		c.Data["json"] = res
		c.ServeJSON()
	}
}
