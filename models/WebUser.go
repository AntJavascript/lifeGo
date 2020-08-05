package models

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type WebUserController struct {
	beego.Controller
}

type EditUser struct {
	UserName string
	FacePath string
}

func (c *WebUserController) GetList() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 获取limit ，每页需要展示的数量
	var pasize = c.GetString("limit")
	limit, err := strconv.ParseInt(pasize, 10, 64)
	// 当前处于第几页
	var page = c.GetString("page")
	currentPage, err := strconv.ParseInt(page, 10, 64)

	var list []User
	num, err := o.QueryTable("User").Limit(limit, currentPage-1).All(&list)
	for index, _ := range list {
		list[index].Password = ""
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

// 获取个人信息
func (c *WebUserController) GetUserInfo() {
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
	// 获取用户信息

	var user_data User
	user_data.Id = flagUser
	err := o.Read(&user_data, "Id")
	user_data.Password = ""

	if err != nil {
		res["data"] = ""
		res["code"] = 401
		res["msg"] = err
	} else {
		res["data"] = user_data
		res["code"] = 200
	}
	c.Data["json"] = res
	c.ServeJSON()
}

// 修改个人信息
func (c *WebUserController) EditUserInfo() {
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

	// 获取用户信息
	var user_data User
	user_data.Id = flagUser
	err := o.Read(&user_data, "Id")

	// 构造数据查询条件
	var params EditUser
	//json数据封装到params对象中
	paramsErr := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	fmt.Println(paramsErr, params)

	var user User // 个人信息
	user.Face = params.FacePath
	user.UserName = params.UserName
	user.Id = flagUser

	_, updateErr := o.Update(&user, "face", "user_name")

	if updateErr != nil {
		res["data"] = ""
		res["code"] = 401
		res["msg"] = err
	} else {
		res["msg"] = ""
		res["data"] = ""
		res["code"] = 200
	}
	c.Data["json"] = res
	c.ServeJSON()
}
