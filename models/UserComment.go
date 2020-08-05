package models

import (
	"encoding/json"
	"fmt"

	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 评论数据表结构
type UserComment struct {
	Id         int
	Comment    string    // 评论内容
	CreateTime time.Time // 评论时间
	Tag        string    // 评论标签
	CommentImg string    // 评论图片
	Star       int       // 评论星级
	OrderId    string    // 订单id
	UserId     int       // 用户id
	ProductId  int       // 产品id
}

type UserCommentController struct {
	beego.Controller
}

// 获取评论列表
func (c *UserCommentController) GetList() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 获取limit ，每页需要展示的数量
	var pasize = c.GetString("limit")
	limit, _ := strconv.ParseInt(pasize, 10, 64)
	// 当前处于第几页
	var page = c.GetString("page")
	currentPage, _ := strconv.ParseInt(page, 10, 64)

	var list []UserComment

	_, err := o.QueryTable("UserComment").OrderBy("-create_time").Limit(limit, (currentPage-1)*limit).All(&list)
	// 计算订单总数
	total, _ := o.QueryTable("UserComment").Count()

	if err != nil {
		res["data"] = ""
		res["code"] = 401
		res["msg"] = err
	} else {
		res["data"] = list
		res["total"] = total
		res["code"] = 200
	}
	c.Data["json"] = res
	c.ServeJSON()
}

// 根据产品id获取评论列表
func (c *UserCommentController) FindProductIdGetComment() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 获取ajax提交的参数
	paramsId := c.Ctx.Input.Param(":id")

	var commentList []UserComment // 评论列表
	var resultList []interface{}  // 返回的json数据

	_, err := o.QueryTable("UserComment").Filter("ProductId", paramsId).OrderBy("-create_time").All(&commentList)

	for _, item := range commentList {
		var returnCOmmentfield = make(map[string]interface{})
		returnCOmmentfield["Comment"] = item.Comment
		returnCOmmentfield["CommentImg"] = item.CommentImg
		returnCOmmentfield["CreateTime"] = item.CreateTime
		returnCOmmentfield["Id"] = item.Id
		returnCOmmentfield["OrderId"] = item.OrderId
		returnCOmmentfield["ProductId"] = item.ProductId
		returnCOmmentfield["Star"] = item.Star
		returnCOmmentfield["UserId"] = item.UserId

		var userInfo User
		_, usererr := o.QueryTable("User").Filter("Id", item.UserId).All(&userInfo)
		fmt.Println(usererr)

		userInfo.Password = "" // 密码字段置为空

		returnCOmmentfield["UserInfo"] = userInfo
		resultList = append(resultList, returnCOmmentfield)
		fmt.Println("评论数据有问题啊", resultList)
	}

	if err != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res["data"] = resultList
		res["code"] = 200
		res["msg"] = "ok"
		c.Data["json"] = res
		c.ServeJSON()
	}

}

// 发布评论
func (c *UserCommentController) SubmitComment() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 构造数据查询条件
	var params UserComment
	//json数据封装到params对象中
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)

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

	params.CreateTime = time.Now()
	params.UserId = flagUser

	var orderInfo UserOrder // 当前订单信息
	orderInfo.OrderId = params.OrderId
	orderInfoErr := o.Read(&orderInfo, "OrderId")

	if orderInfoErr != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = orderInfoErr
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	fmt.Println("当前订单信息状态", orderInfo.OrderStatus)
	if orderInfo.OrderStatus == 6 {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = "订单已经评价过了"
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if orderInfo.OrderStatus != 5 {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = "当前订单不可评价"
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	_, Inserterr := o.Insert(&params)

	if Inserterr != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		// 改变订单状态
		var User_order UserOrder
		User_order.OrderStatus = 6 // 已评价
		User_order.OrderId = params.OrderId
		orderId, UserOrdererr := o.QueryTable("UserOrder").Filter("OrderId", params.OrderId).Update(orm.Params{
			"OrderStatus": User_order.OrderStatus,
		})
		fmt.Println(orderId, UserOrdererr) // 打印

		res["data"] = ""
		res["code"] = 200
		res["msg"] = "ok"
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
	}
}
