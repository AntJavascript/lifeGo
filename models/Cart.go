package models

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Cart struct {
	Id         int
	ProductId  int // 商品id
	CartNumber int // 购物车数量
	Attr       int // 属性id
	UserId     int // 用户id
	Status     int // 购物车状态
	IsBuy      int // 是否已被购买
}
type CartJson struct {
	Cart         Cart
	AttrList     ProductSpecs
	ProductImg   string
	ProductTitle string
	ProductPrice float64
	ProductStock int
}
type CartController struct {
	beego.Controller
}

// 修改购物车状态的参数
type CheckStatus struct {
	CartId  int
	IsCheck int
}

// 要删除的商品id集合
type CartIds struct {
	Ids []int
}

// 获取购物车列表
func (c *CartController) FindList() {
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

	var cartList []Cart
	num, err := o.QueryTable("Cart").Filter("user_id", flagUser).Filter("is_buy", "0").All(&cartList)
	if err != nil {
		res["code"] = 401
	} else {
		res["total"] = num
	}
	CartJsonData := make([]CartJson, 0)
	for i, v := range cartList {
		var list []ProductSpecs
		var jsonCart CartJson
		_, err := o.QueryTable("ProductSpecs").Filter("ProductId", v.ProductId).All(&list)
		fmt.Println(err)
		for index, item := range list {
			// 如果ProductSpecs.Id字段和Cart.Id一样
			if item.Id == cartList[i].Attr {
				jsonCart.AttrList = list[index]
			}
		}
		var product Product

		_, productErr := o.QueryTable("Product").Filter("Id", v.ProductId).All(&product)
		fmt.Println(productErr)

		jsonCart.ProductImg = product.Thumbnail
		jsonCart.ProductTitle = product.Title
		jsonCart.ProductPrice = product.Price
		jsonCart.ProductStock = product.Stock

		jsonCart.Cart = cartList[i]
		CartJsonData = append(CartJsonData, jsonCart)
	}
	res["code"] = 200
	res["msg"] = err
	res["data"] = CartJsonData
	c.Data["json"] = res
	c.ServeJSON()
}

// 添加购物车
func (c *CartController) AddCart() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	var loginSession = c.GetSession("LoginSession")
	var loginSessionUser User
	// 判断类型是否是User类型
	if v, ok := loginSession.(User); ok {
		loginSessionUser = v
	}
	if loginSession == nil {
		res["code"] = 401
		res["msg"] = "请先登录"
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	var params Cart
	paramsErr := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if paramsErr != nil {
		res["code"] = 400
		res["msg"] = paramsErr
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	params.UserId = loginSessionUser.Id
	id, err := o.Insert(&params)
	if err != nil {
		res["code"] = 401
		res["msg"] = "添加购物车失败"
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	res["code"] = 200
	res["msg"] = "添加购物车成功"
	res["data"] = id
	c.Data["json"] = res
	c.ServeJSON()
}

// 修改购物车状态
func (c *CartController) CartStatus() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	var params CheckStatus
	// 解析参数
	paramsErr := json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	if paramsErr != nil {
		res["code"] = 400
		res["msg"] = paramsErr
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	var sqlData Cart
	sqlData.Id = params.CartId
	sqlData.Status = params.IsCheck
	fmt.Println(sqlData)
	_, err := o.Update(&sqlData, "Status")
	if err != nil {
		res["code"] = 401
		res["msg"] = err
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	res["code"] = 200
	res["msg"] = ""
	c.Data["json"] = res
	c.ServeJSON()
}

// 获取已选中购物车列表
func (c *CartController) FindSelectList() {
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

	var cartList []Cart
	num, err := o.QueryTable("Cart").Filter("user_id", flagUser).Filter("Status", 1).Filter("is_buy", 0).All(&cartList)
	if err != nil {
		res["code"] = 401
		res["msg"] = err
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	} else {
		res["total"] = num
	}
	CartJsonData := make([]CartJson, 0)
	for i, v := range cartList {
		var list []ProductSpecs
		var jsonCart CartJson
		_, err := o.QueryTable("ProductSpecs").Filter("ProductId", v.ProductId).All(&list)
		fmt.Println(err)
		for index, item := range list {
			// 如果ProductSpecs.Id字段和Cart.Id一样
			if item.Id == cartList[i].Attr {
				jsonCart.AttrList = list[index]
			}
		}
		var product Product
		_, productErr := o.QueryTable("Product").Filter("Id", v.ProductId).All(&product)
		fmt.Println(productErr)

		jsonCart.ProductImg = product.Thumbnail
		jsonCart.ProductTitle = product.Title
		jsonCart.ProductPrice = product.Price
		jsonCart.ProductStock = product.Stock

		jsonCart.Cart = cartList[i]
		CartJsonData = append(CartJsonData, jsonCart)
	}
	res["code"] = 200
	res["msg"] = err
	res["data"] = CartJsonData
	c.Data["json"] = res
	c.ServeJSON()
}

// 修改购物车数量
func (c *CartController) EditCartNumber() {
	res := make(map[string]interface{})
	o := orm.NewOrm()

	// 构造数据查询条件
	var params Cart
	//json数据封装到params对象中
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	fmt.Println(err)

	// 更新购物车数量
	_, Updateerr := o.Update(&params, "CartNumber")

	if Updateerr == nil {
		res["code"] = 200
		res["msg"] = "ok"
		res["data"] = ""
		c.Data["json"] = res
	} else {
		res["code"] = 400
		res["msg"] = Updateerr
		res["data"] = ""
		c.Data["json"] = res
	}
	c.ServeJSON()
}

// 删除购物车
func (c *CartController) DelCart() {
	res := make(map[string]interface{})
	o := orm.NewOrm()

	// 构造数据查询条件
	var params CartIds
	//json数据封装到params对象中
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	fmt.Println(err, params)

	for _, value := range params.Ids {
		// 删除购物车
		_, sqlErr := o.QueryTable("Cart").Filter("id", value).Delete()
		fmt.Println(sqlErr)
	}

	res["code"] = 200
	res["msg"] = "ok"
	res["data"] = ""
	c.Data["json"] = res
	c.ServeJSON()
}
