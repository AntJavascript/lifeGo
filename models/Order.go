package models

import (
	"encoding/json"
	"fmt"

	//	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserOrder struct {
	Id              int       // id
	ProductId       string    // 产品id，可以多个逗号分隔
	CartId          string    // 购物车id，可以多个逗号分隔
	RealPayAmount   float64   // 实际支付金额
	CouponId        int       // 优惠券id
	CouponAmount    float64   // 优惠券金额
	Phone           string    // 收货人电话
	Name            string    // 收货姓名
	Address         string    // 收货地址
	PayType         int       // 支付方式
	FreightAmount   int       // 运费
	Remarks         string    // 备注
	OrderStatus     int       // 订单状态  1、待付款 2、待发货 3、配送中 4、待收货 5、待评价 6、已评价 7、已取消
	UserId          int       // 用户id
	Logistics       string    // 物流名称
	LogisticsNumber string    // 物流单号
	BuyNumber       int       // 购买数量
	OrderId         string    // 订单id
	CreateTime      time.Time // 创建时间
	TotalPrice      float64   // 商品总额
	SendTime        string    // 发货时间
	PayTime         string    // 付款时间
	ReceiveTime     string    // 确认收货时间
	AttrId          string    // 商品属性id
}

// 个人订单列表
type UserOrderList struct {
	Order            UserOrder      // 订单信息
	ProductList      []Product      // 商品清单
	CartList         []Cart         // 购物车
	ProductSpecsList []ProductSpecs // 商品属性组合
	Comment          []UserComment  // 评论信息
}

// 支付订单参数
type payParams struct {
	Id          int //  数据表id
	PayType     int // 支付方式
	OrderStatus int // 订单状态
}

// 取消订单参数
type CancelOrderParams struct {
	OrderId     string // 订单id
	Id          int
	OrderStatus int
}

// 发货参数
type sendParams struct {
	Id              int
	OrderId         string // 订单id
	Logistics       string // 物流名称
	LogisticsNumber string // 物流单号
	SendTime        string // 发货时间
	OrderStatus     int
}

// 收货参数
type ReceiveParams struct {
	Id          int
	OrderId     string // 订单id
	ReceiveTime string // 发货时间
}

type OrderController struct {
	beego.Controller
}

// 设置购物车商品状态
func ChangeCartStatus(cartId string) {
	CartIdArr := strings.Split(cartId, ",")
	o := orm.NewOrm()
	index := strings.Index(cartId, ",")

	if index == -1 {
		formatId, err := strconv.Atoi(cartId)
		if err == nil {
			var cart Cart
			cart.Id = formatId
			cart.IsBuy = 1
			o.Update(&cart, "IsBuy")
		}
	} else {
		for _, id := range CartIdArr {
			formatId, err := strconv.Atoi(id)
			if err == nil {
				var cart Cart
				cart.Id = formatId
				cart.IsBuy = 1
				o.Update(&cart, "IsBuy")
			}
		}
	}
}

// 设置优惠券状态
func ChangeCouponStatus(couponId int) {
	o := orm.NewOrm()

	var coupon UserCoupon
	coupon.Id = couponId
	coupon.IsUse = 1
	o.Update(&coupon, "IsUse")
}

// 生成订单
func (c *OrderController) Buy() {
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
	var params UserOrder
	//json数据封装到params对象中
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	params.UserId = flagUser
	id := strconv.FormatInt(rand.Int63()+time.Now().Unix(), 10) // 生成订单id
	params.OrderId = id
	params.OrderStatus = 1 // 设置订单状态
	params.CreateTime = time.Now()

	if err != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
		return
	} else {
		// 先查询一下订单id是否已存在
		orderId, err := o.Insert(&params)
		if err != nil {
			res["data"] = ""
			res["code"] = 401
			res["msg"] = err
			// 返回json格式给前端界面
			c.Data["json"] = res
			c.ServeJSON()
			return
		}
		// 设置购物车商品状态
		ChangeCartStatus(params.CartId)

		// 设置优惠券状态
		ChangeCouponStatus(params.CouponId)

		res["data"] = orderId
		res["code"] = 200
		res["msg"] = err
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
	}
}

// 后台展示订单列表
func (c *OrderController) GetList() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 获取limit ，每页需要展示的数量
	var pasize = c.GetString("limit")
	limit, err := strconv.ParseInt(pasize, 10, 64)
	// 当前处于第几页
	var page = c.GetString("page")
	currentPage, err := strconv.ParseInt(page, 10, 64)

	var list []UserOrder
	var lists []UserOrderList           //个人订单列表
	var UserOrderListJSON UserOrderList //订单JSON

	num, err := o.QueryTable("UserOrder").OrderBy("-create_time").Limit(limit, (currentPage-1)*limit).All(&list)
	fmt.Println(num)
	// 计算订单总数
	total, _ := o.QueryTable("UserOrder").Count()

	for _, value := range list {
		var products []Product
		// 判断是否有多个商品id
		indexs := strings.Index(value.ProductId, ",") // 已逗号分隔商品id
		if indexs == -1 {
			// 如果只有一个商品就请求当前商品的数据
			split := value.ProductId
			var product Product
			_, productErr := o.QueryTable("Product").Filter("Id", split).All(&product)

			// 如果没错误就添加进数组
			if productErr == nil {
				products = append(products, product)
			}
		} else {
			split := strings.Split(value.ProductId, ",") // 已逗号分隔商品id
			// 如果存在多个商品就循环请求商品信息
			for _, productID := range split {
				var product Product
				_, productErr := o.QueryTable("Product").Filter("Id", productID).All(&product)

				// 如果没错误就添加进数组
				if productErr == nil {
					products = append(products, product)
				}
			}
		}

		var carts []Cart
		// 判断是否有多个购物车id
		cartIndex := strings.Index(value.CartId, ",") // 已逗号分隔商品id
		if cartIndex == -1 {
			// 如果只有一个商品就请求当前商品的数据
			split := value.CartId
			var cart Cart
			_, CartErr := o.QueryTable("Cart").Filter("Id", split).All(&cart)

			// 如果没错误就添加进数组
			if CartErr == nil {
				carts = append(carts, cart)
			}
		} else {
			split := strings.Split(value.CartId, ",") // 已逗号分隔商品id
			// 如果存在多个商品就循环请求商品信息
			for _, CartID := range split {
				var cart Cart
				_, CartErr := o.QueryTable("Cart").Filter("Id", CartID).All(&cart)

				// 如果没错误就添加进数组
				if CartErr == nil {
					carts = append(carts, cart)
				}
			}
		}
		var productSpecs []ProductSpecs
		// 查询商品属性组合
		for _, ProductSpecsItem := range carts {
			var productSpec ProductSpecs
			_, ProductSpecErr := o.QueryTable("ProductSpecs").Filter("Id", ProductSpecsItem.Attr).All(&productSpec)

			// 如果没错误就添加进数组
			if ProductSpecErr == nil {
				productSpecs = append(productSpecs, productSpec)
			}
		}

		UserOrderListJSON.Order = value
		UserOrderListJSON.ProductList = products
		UserOrderListJSON.CartList = carts
		UserOrderListJSON.ProductSpecsList = productSpecs

		lists = append(lists, UserOrderListJSON)

	}

	if err != nil {
		res["data"] = ""
		res["code"] = 401
		res["msg"] = err
	} else {
		res["data"] = lists
		res["total"] = total
		res["code"] = 200
	}
	c.Data["json"] = res
	c.ServeJSON()
}

// 个人订单列表
func (c *OrderController) GetUserList() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	var Type = c.GetString("type") // 订单类型 0、全部订单 1、待付款 2、待发货 3、待收货 4、待评价 5、已完成
	orderType, _ := strconv.Atoi(Type)
	// 获取limit ，每页需要展示的数量
	var pasize = c.GetString("limit")
	limit, _ := strconv.ParseInt(pasize, 10, 64)
	// 当前处于第几页
	var page = c.GetString("page")
	currentPage, err := strconv.ParseInt(page, 10, 64)
	// 获取登录人session
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

	var lists []UserOrderList           //个人订单列表
	var UserOrderListJSON UserOrderList //订单JSON

	var OrderList []UserOrder
	var num int64
	var orderErr error
	if orderType == 0 {
		// 全部订单
		n, e := o.QueryTable("UserOrder").Filter("user_id", flagUser).OrderBy("-create_time").Limit(limit, (currentPage-1)*10).All(&OrderList)
		num = n
		orderErr = e
	} else {
		n, e := o.QueryTable("UserOrder").Filter("user_id", flagUser).OrderBy("-create_time").Filter("OrderStatus", orderType).Limit(limit, (currentPage-1)*10).All(&OrderList)
		num = n
		orderErr = e
	}

	for _, value := range OrderList {
		// 订单状态评价信息
		var commentlists []UserComment

		o.QueryTable("UserComment").Filter("OrderId", value.OrderId).All(&commentlists)

		var products []Product
		// 判断是否有多个商品id
		indexs := strings.Index(value.ProductId, ",") // 已逗号分隔商品id
		if indexs == -1 {
			// 如果只有一个商品就请求当前商品的数据
			split := value.ProductId
			var product Product
			_, productErr := o.QueryTable("Product").Filter("Id", split).All(&product)

			// 如果没错误就添加进数组
			if productErr == nil {
				products = append(products, product)
			}
		} else {
			split := strings.Split(value.ProductId, ",") // 已逗号分隔商品id
			// 如果存在多个商品就循环请求商品信息
			for _, productID := range split {
				var product Product
				_, productErr := o.QueryTable("Product").Filter("Id", productID).All(&product)

				// 如果没错误就添加进数组
				if productErr == nil {
					products = append(products, product)
				}
			}
		}

		var carts []Cart
		// 判断是否有多个购物车id
		cartIndex := strings.Index(value.CartId, ",") // 已逗号分隔商品id
		if cartIndex == -1 {
			// 如果只有一个商品就请求当前商品的数据
			split := value.CartId
			var cart Cart
			_, CartErr := o.QueryTable("Cart").Filter("Id", split).All(&cart)

			// 如果没错误就添加进数组
			if CartErr == nil {
				carts = append(carts, cart)
			}
			fmt.Println("购物车id", split)
		} else {
			split := strings.Split(value.CartId, ",") // 已逗号分隔商品id
			// 如果存在多个商品就循环请求商品信息
			for _, CartID := range split {
				var cart Cart
				_, CartErr := o.QueryTable("Cart").Filter("Id", CartID).All(&cart)

				// 如果没错误就添加进数组
				if CartErr == nil {
					carts = append(carts, cart)
				}
			}
		}
		var productSpecs []ProductSpecs
		// 查询商品属性组合
		for _, ProductSpecsItem := range carts {
			var productSpec ProductSpecs
			_, ProductSpecErr := o.QueryTable("ProductSpecs").Filter("Id", ProductSpecsItem.Attr).All(&productSpec)

			// 如果没错误就添加进数组
			if ProductSpecErr == nil {
				productSpecs = append(productSpecs, productSpec)
			}
		}

		UserOrderListJSON.Order = value
		UserOrderListJSON.ProductList = products
		UserOrderListJSON.CartList = carts
		UserOrderListJSON.ProductSpecsList = productSpecs
		UserOrderListJSON.Comment = commentlists

		lists = append(lists, UserOrderListJSON)

	}

	if orderErr != nil {
		res["data"] = ""
		res["code"] = 401
		res["msg"] = err
	} else {
		res["data"] = lists
		res["total"] = num
		res["code"] = 200
	}
	c.Data["json"] = res
	c.ServeJSON()
}

// 获取订单详情
func (c *OrderController) GetUserOrderDetail() {
	res := make(map[string]interface{})
	o := orm.NewOrm()

	// 获取ajax提交的参数
	paramsId := c.Ctx.Input.Param(":id")

	// 获取登录人session
	loginSession := c.GetSession("LoginSession")
	if loginSession == nil {
		res["code"] = 401
		res["msg"] = "请先登录"
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 登录人标识
	var flagUser int
	if v, ok := loginSession.(User); ok {
		flagUser = v.Id
	}
	id, _ := strconv.Atoi(paramsId)

	var UserOrderListJSON UserOrderList //订单JSON

	var OrderInfo UserOrder
	OrderInfo.Id = id           // 订单id
	OrderInfo.UserId = flagUser // 登录人id
	_, orderErr := o.QueryTable("UserOrder").Filter("user_id", flagUser).Filter("id", OrderInfo.Id).All(&OrderInfo)

	var products []Product
	// 判断是否有多个商品id
	indexs := strings.Index(OrderInfo.ProductId, ",") // 已逗号分隔商品id
	if indexs == -1 {
		// 如果只有一个商品就请求当前商品的数据
		split := OrderInfo.ProductId
		var product Product
		_, productErr := o.QueryTable("Product").Filter("Id", split).All(&product)

		// 如果没错误就添加进数组
		if productErr == nil {
			products = append(products, product)
		}
	} else {
		split := strings.Split(OrderInfo.ProductId, ",") // 已逗号分隔商品id
		// 如果存在多个商品就循环请求商品信息
		for _, productID := range split {
			var product Product
			_, productErr := o.QueryTable("Product").Filter("Id", productID).All(&product)

			// 如果没错误就添加进数组
			if productErr == nil {
				products = append(products, product)
			}
		}
	}

	var carts []Cart
	// 判断是否有多个购物车id
	cartIndex := strings.Index(OrderInfo.CartId, ",") // 已逗号分隔商品id
	if cartIndex == -1 {
		// 如果只有一个商品就请求当前商品的数据
		split := OrderInfo.CartId
		var cart Cart
		_, CartErr := o.QueryTable("Cart").Filter("Id", split).All(&cart)

		// 如果没错误就添加进数组
		if CartErr == nil {
			carts = append(carts, cart)
		}
	} else {
		split := strings.Split(OrderInfo.CartId, ",") // 已逗号分隔商品id
		// 如果存在多个商品就循环请求商品信息
		for _, CartID := range split {
			var cart Cart
			_, CartErr := o.QueryTable("Cart").Filter("Id", CartID).All(&cart)

			// 如果没错误就添加进数组
			if CartErr == nil {
				carts = append(carts, cart)
			}
		}
	}
	var productSpecs []ProductSpecs
	// 查询商品属性组合
	for _, ProductSpecsItem := range carts {
		var productSpec ProductSpecs
		_, ProductSpecErr := o.QueryTable("ProductSpecs").Filter("Id", ProductSpecsItem.Attr).All(&productSpec)

		// 如果没错误就添加进数组
		if ProductSpecErr == nil {
			productSpecs = append(productSpecs, productSpec)
		}
	}

	UserOrderListJSON.Order = OrderInfo
	UserOrderListJSON.ProductList = products
	UserOrderListJSON.CartList = carts
	UserOrderListJSON.ProductSpecsList = productSpecs

	if orderErr != nil {
		res["data"] = ""
		res["code"] = 401
		res["msg"] = orderErr
	} else {
		res["data"] = UserOrderListJSON
		res["code"] = 200
		res["msg"] = ""
	}
	c.Data["json"] = res
	c.ServeJSON()
}

// 后台获取订单详情
func (c *OrderController) AdminGetUserOrderDetail() {
	res := make(map[string]interface{})
	o := orm.NewOrm()

	// 获取ajax提交的参数
	paramsId := c.Ctx.Input.Param(":id")

	id, _ := strconv.Atoi(paramsId)

	var UserOrderListJSON UserOrderList //订单JSON

	var OrderInfo UserOrder
	OrderInfo.Id = id // 订单id
	_, orderErr := o.QueryTable("UserOrder").Filter("id", OrderInfo.Id).All(&OrderInfo)

	var products []Product
	// 判断是否有多个商品id
	indexs := strings.Index(OrderInfo.ProductId, ",") // 已逗号分隔商品id
	if indexs == -1 {
		// 如果只有一个商品就请求当前商品的数据
		split := OrderInfo.ProductId
		var product Product
		_, productErr := o.QueryTable("Product").Filter("Id", split).All(&product)

		// 如果没错误就添加进数组
		if productErr == nil {
			products = append(products, product)
		}
	} else {
		split := strings.Split(OrderInfo.ProductId, ",") // 已逗号分隔商品id
		// 如果存在多个商品就循环请求商品信息
		for _, productID := range split {
			var product Product
			_, productErr := o.QueryTable("Product").Filter("Id", productID).All(&product)

			// 如果没错误就添加进数组
			if productErr == nil {
				products = append(products, product)
			}
		}
	}

	var carts []Cart
	// 判断是否有多个购物车id
	cartIndex := strings.Index(OrderInfo.CartId, ",") // 已逗号分隔商品id
	if cartIndex == -1 {
		// 如果只有一个商品就请求当前商品的数据
		split := OrderInfo.CartId
		var cart Cart
		_, CartErr := o.QueryTable("Cart").Filter("Id", split).All(&cart)

		// 如果没错误就添加进数组
		if CartErr == nil {
			carts = append(carts, cart)
		}
	} else {
		split := strings.Split(OrderInfo.CartId, ",") // 已逗号分隔商品id
		// 如果存在多个商品就循环请求商品信息
		for _, CartID := range split {
			var cart Cart
			_, CartErr := o.QueryTable("Cart").Filter("Id", CartID).All(&cart)

			// 如果没错误就添加进数组
			if CartErr == nil {
				carts = append(carts, cart)
			}
		}
	}
	var productSpecs []ProductSpecs
	// 查询商品属性组合
	for _, ProductSpecsItem := range carts {
		var productSpec ProductSpecs
		_, ProductSpecErr := o.QueryTable("ProductSpecs").Filter("Id", ProductSpecsItem.Attr).All(&productSpec)

		// 如果没错误就添加进数组
		if ProductSpecErr == nil {
			productSpecs = append(productSpecs, productSpec)
		}
	}

	UserOrderListJSON.Order = OrderInfo
	UserOrderListJSON.ProductList = products
	UserOrderListJSON.CartList = carts
	UserOrderListJSON.ProductSpecsList = productSpecs

	if orderErr != nil {
		res["data"] = ""
		res["code"] = 401
		res["msg"] = orderErr
	} else {
		res["data"] = UserOrderListJSON
		res["code"] = 200
		res["msg"] = ""
	}
	c.Data["json"] = res
	c.ServeJSON()
}

// 支付订单
func (c *OrderController) Pay() {
	res := make(map[string]interface{})
	o := orm.NewOrm()

	// 获取登录人session
	loginSession := c.GetSession("LoginSession")
	if loginSession == nil {
		res["code"] = 401
		res["msg"] = "请先登录"
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 构造数据查询条件
	var params payParams
	//json数据封装到params对象中
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	var sqlParams UserOrder

	sqlParams.Id = params.Id
	sqlParams.PayType = params.PayType
	sqlParams.OrderStatus = 2 // 待发货

	// 先查询订单状态是否是待支付状态
	var findUserOrder UserOrder
	o.QueryTable("UserOrder").Filter("Id", params.Id).All(&findUserOrder)

	// 如果订单状态不等于1
	if findUserOrder.OrderStatus != 1 {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = "订单不处于待支付状态"
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	indexs := strings.Index(findUserOrder.AttrId, ",") // 已逗号分隔商品id
	fmt.Println("已逗号分隔商品id", indexs)

	// 订单存在多个商品
	split := strings.Split(findUserOrder.AttrId, ",")           // 用逗号分隔商品属性id
	cartSplit := strings.Split(findUserOrder.CartId, ",")       // 用逗号分隔购物车id
	productSplit := strings.Split(findUserOrder.ProductId, ",") // 用逗号分隔商品id

	// 商品存在组合属性
	if findUserOrder.AttrId != "0" {

		// 如果存在多个商品就循环请求商品信息
		for index, attrID := range split {

			var cart_data Cart // 购物车数据
			_, carterr := o.QueryTable("cart").Filter("Id", cartSplit[index]).All(&cart_data)
			fmt.Println(carterr)

			// 如果商品属性存在（!=0）
			if attrID != "0" {
				// 获取商品属性组合
				var productSpecs ProductSpecs
				_, attrErr := o.QueryTable("ProductSpecs").Filter("Id", attrID).All(&productSpecs)
				fmt.Println(attrErr)
				intId, _ := strconv.Atoi(productSpecs.ProductStock) // 库存

				// 如果购买数量大于剩余库存，则直接返回
				if cart_data.CartNumber > intId {
					res["data"] = ""
					res["code"] = 400
					res["msg"] = "库存不足"
					c.Data["json"] = res
					c.ServeJSON()
					return
				}

				// 减少库存操作
				productSpecs.ProductStock = strconv.Itoa(intId - cart_data.CartNumber)
				_, Updateerr := o.Update(&productSpecs, "ProductStock")
				fmt.Println(Updateerr)
			} else {
				// 减少添加商品时的库存
				var product_data Product
				_, attrErr := o.QueryTable("Product").Filter("Id", productSplit[index]).All(&product_data)
				fmt.Println(attrErr)
				// 如果购买数量大于剩余库存，则直接返回
				if cart_data.CartNumber > product_data.Stock {
					res["data"] = ""
					res["code"] = 400
					res["msg"] = "库存不足"
					c.Data["json"] = res
					c.ServeJSON()
					return
				}
				// 减少库存操作
				product_data.Stock = product_data.Stock - cart_data.CartNumber
				_, Updateerr := o.Update(&product_data, "Stock")
				fmt.Println(Updateerr)
			}

		}

	} else {

		var cart_data Cart // 购物车数据
		_, carterr := o.QueryTable("cart").Filter("Id", findUserOrder.ProductId).All(&cart_data)
		fmt.Println("购物车数据", carterr)

		// 产品表库存字段 ProductId
		var product Product
		productId, _ := strconv.Atoi(findUserOrder.ProductId)
		o.QueryTable("product").Filter("Id", productId).All(&product)

		if product.Stock >= cart_data.CartNumber {
			product.Stock = product.Stock - cart_data.CartNumber
			_, Updateerr := o.Update(&product, "Stock")
			fmt.Println(Updateerr)
		} else {
			res["data"] = ""
			res["code"] = 400
			res["msg"] = "库存不足"
			c.Data["json"] = res
			c.ServeJSON()
			return
		}
	}

	// 产品表库存字段   ProductId
	var product_data Product
	o.QueryTable("product").Filter("Id", findUserOrder.ProductId).All(&product_data)

	for _, productId := range productSplit {
		var product_data Product
		o.QueryTable("product").Filter("Id", productId).All(&product_data)

		// 产品销量加1
		product_data.SalesVolume += 1

		o.Update(&product_data, "SalesVolume")
	}

	sqlParams.PayTime = strconv.FormatInt(time.Now().Unix(), 10) // 付款时间

	_, updetaErr := o.Update(&sqlParams, "PayType", "OrderStatus", "PayTime")

	if updetaErr != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res["data"] = ""
		res["code"] = 200
		res["msg"] = "ok"
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
	}
}

// 取消订单
func (c *OrderController) CancelOrder() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 获取登录人session
	loginSession := c.GetSession("LoginSession")
	if loginSession == nil {
		res["code"] = 401
		res["msg"] = "请先登录"
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 构造数据查询条件
	var params CancelOrderParams
	//json数据封装到params对象中
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	if err != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	var sqlParams UserOrder

	sqlParams.OrderId = params.OrderId
	sqlParams.Id = params.Id
	sqlParams.OrderStatus = params.OrderStatus

	_, updetaErr := o.Update(&sqlParams, "OrderStatus")
	if updetaErr != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res["data"] = ""
		res["code"] = 200
		res["msg"] = "ok"
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
	}
}

// 发货处理
func (c *OrderController) SendOrder() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 构造数据查询条件
	var params sendParams
	//json数据封装到params对象中
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	if err != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if params.Id <= 0 {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = "id不能为空"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if params.OrderId == "" {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = "订单id不能为空"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if params.Logistics == "" {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = "物流不能为空"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if params.LogisticsNumber == "" {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = "物流单号不能为空"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var sqlParams UserOrder

	sqlParams.OrderId = params.OrderId
	sqlParams.Id = params.Id
	sqlParams.Logistics = params.Logistics                        //物流公司
	sqlParams.LogisticsNumber = params.LogisticsNumber            // 物流单号
	sqlParams.SendTime = strconv.FormatInt(time.Now().Unix(), 10) // 发货时间
	sqlParams.OrderStatus = 3

	_, updetaErr := o.Update(&sqlParams, "Logistics", "LogisticsNumber", "SendTime", "OrderStatus")

	if updetaErr != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res["data"] = ""
		res["code"] = 200
		res["msg"] = "ok"
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
	}
}

//确认收货
func (c *OrderController) Receiving() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 构造数据查询条件
	var params ReceiveParams
	//json数据封装到params对象中
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	if err != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if params.Id <= 0 {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = "id不能为空"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if params.OrderId == "" {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = "订单id不能为空"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	var sqlParams UserOrder

	sqlParams.OrderId = params.OrderId
	sqlParams.Id = params.Id
	sqlParams.ReceiveTime = strconv.FormatInt(time.Now().Unix(), 10) // 确认收货时间
	sqlParams.OrderStatus = 5

	_, updetaErr := o.Update(&sqlParams, "ReceiveTime", "OrderStatus")

	if updetaErr != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res["data"] = ""
		res["code"] = 200
		res["msg"] = "ok"
		// 返回json格式给前端界面
		c.Data["json"] = res
		c.ServeJSON()
	}
}
