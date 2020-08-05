package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 商品表字段
type Product struct {
	Id              int       // id
	Title           string    // 名称
	ProductDesc     string    // 描述
	Price           float64   // 价格
	Detail          string    // 详情
	AttributeList   string    // 属性
	Stock           int       // 库存
	Status          int       // 状态
	OrignPrice      float64   // 市场价
	CategoryId      int       // 分类
	CategoryName    string    // 分类名称
	CategorySubId   int       // 二级分类
	CategorySubName string    // 二级分类名称
	CreateTime      time.Time // 创建时间
	UpdateTime      time.Time // 修改时间
	Freight         int       // 运费
	PlaceOrigin     string    // 产地
	Album           string    // 相册
	ProductSpecs    string    // 组合属性
	Thumbnail       string    // 缩略图
	IsDel           int       // 是否删除
	SalesVolume     int       // 销量
	Hot             int       // 是否热销
	Recommend       int       // 是否热销推荐
}

// 商品属性字段
type ProductSpecs struct {
	Id           int       // id
	ProductId    int       // 商品id
	Specs        string    // 属性
	CreateTime   time.Time // 创建时间
	ProductStock string    // 库存
	ProductPrice string    // 价格
}

// 返回前端字段适配flutter
type ResultData struct {
	Id            int     // id
	Title         string  // 名称
	ProductDesc   string  // 描述
	Price         float64 // 价格
	Detail        string  // 详情
	AttributeList string  // 属性
	Stock         int     // 库存
	Status        int     // 状态
	OrignPrice    float64 // 市场价
	Freight       int     // 运费
	PlaceOrigin   string  // 产地
	Album         string  // 相册
	Thumbnail     string  // 缩略图
	SalesVolume   int     // 销量
	Hot           int     // 是否热销
	Recommend     int     // 是否热销推荐
	SpecsList     []interface{}
}

// 删除商品参数
type delProductParams struct {
	Id int // 商品id
}

type ProductController struct {
	beego.Controller
}

// 获取商品列表
func (c *ProductController) GetList() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 获取limit ，每页需要展示的数量
	var pasize = c.GetString("limit")
	limit, _ := strconv.ParseInt(pasize, 10, 64)
	// 当前处于第几页
	var page = c.GetString("page")
	currentPage, _ := strconv.ParseInt(page, 10, 64)

	// 排序 1、升序；2、降序
	var order = c.GetString("order")
	fmt.Println(order)
	var list []Product
	var json = make([]interface{}, 0)
	_, err := o.QueryTable("Product").OrderBy("-sales_volume").Limit(limit, (currentPage-1)*limit).All(&list)
	total, totalErr := o.QueryTable("Product").Count()
	for _, value := range list {
		json = append(json, value)
	}
	if err != nil {
		//		fmt.Println(num)
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

// 新增商品
func (c *ProductController) AddPost() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 构造数据查询条件
	var params Product
	data := c.Ctx.Input.RequestBody
	//json数据封装到params对象中
	err := json.Unmarshal(data, &params)
	if err != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
	} else {
		// 获取当前时间
		params.CreateTime = time.Now()
		// 插入数据库, 返回id和err错误
		id, err := o.Insert(&params)
		if err != nil {
			res["data"] = id
			res["code"] = 401
			res["msg"] = err
		} else {
			// 构造商品属性字段
			var specs []ProductSpecs
			// 解析参数
			err := json.Unmarshal([]byte(params.ProductSpecs), &specs)
			if err != nil {
				res["data"] = id
				res["code"] = 200
				res["msg"] = "商品插入失败"
			} else {
				// 循环给每个属性组合添加商品id
				for index, _ := range specs {
					specs[index].ProductId = int(id)
				}
				// 批量插入商品属性表，返回数量successNums和错误err
				successNums, err := o.InsertMulti(1000, specs)
				if err != nil {
					res["data"] = successNums
					res["code"] = 401
					res["msg"] = err
				} else {
					fmt.Println(err)
					res["data"] = id
					res["code"] = 200
					res["msg"] = "新增成功"
				}
			}

		}
	}
	// 返回json格式给前端界面
	c.Data["json"] = res
	c.ServeJSON()
}

// 根据id获取数据
func (c *ProductController) FindIdData() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 获取ajax提交的参数
	paramsId := c.Ctx.Input.Param(":id")
	parseId, _ := strconv.Atoi(paramsId)
	params := Product{Id: parseId}

	err := o.Read(&params)
	if err != nil {
		res["code"] = 401
		res["msg"] = err
	} else {
		var attrList []ProductSpecs
		num, err := o.QueryTable("ProductSpecs").Filter("product_id", paramsId).All(&attrList)
		if err != nil {
			res["code"] = 401
		} else {
			res["num"] = num
		}
		// 返回给前端的参数为了兼容flutter
		var productSpecsList ResultData

		productSpecsList.Id = params.Id
		productSpecsList.Title = params.Title
		productSpecsList.Album = params.Album
		productSpecsList.Detail = params.Detail
		productSpecsList.Freight = params.Freight
		productSpecsList.Hot = params.Hot
		productSpecsList.OrignPrice = params.OrignPrice
		productSpecsList.PlaceOrigin = params.PlaceOrigin
		productSpecsList.Price = params.Price
		productSpecsList.ProductDesc = params.ProductDesc
		productSpecsList.Recommend = params.Recommend
		productSpecsList.SalesVolume = params.SalesVolume
		productSpecsList.Stock = params.Stock
		productSpecsList.Thumbnail = params.Thumbnail

		dataArr := []interface{}{}

		dataArrerr := json.Unmarshal([]byte(params.AttributeList), &dataArr)
		productSpecsList.SpecsList = dataArr

		if dataArrerr != nil {
			res["code"] = 401
			res["attrList"] = ""
			res["msg"] = dataArrerr
			res["data"] = ""
		} else {
			res["code"] = 200
			res["attrList"] = attrList
			res["msg"] = ""
			res["data"] = productSpecsList
		}

	}
	// 返回json格式给前端界面
	c.Data["json"] = res
	c.ServeJSON()
}

// 修改商品详情
func (c *ProductController) EditProductDetail() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 构造数据查询条件
	var params Product
	data := c.Ctx.Input.RequestBody
	//json数据封装到params对象中
	err := json.Unmarshal(data, &params)
	if err != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
	} else {
		// 获取当前时间
		params.UpdateTime = time.Now()
		// 插入数据库, 返回id和err错误
		id, err := o.Update(&params)
		if err != nil {
			res["data"] = ""
			res["code"] = 401
			res["msg"] = err
		} else {
			// 构造商品属性字段
			var specs []ProductSpecs
			// 解析参数
			err := json.Unmarshal([]byte(params.ProductSpecs), &specs)
			if err != nil {
				fmt.Println(err)
			} else {
				// 批量插入商品属性表，返回数量successNums和错误err
				for _, item := range specs {
					if item.Id > 0 {
						// 如果存在就更新数据
						successNums, err := o.Update(&item)
						fmt.Println(successNums)
						fmt.Println(err)
					} else {
						// 如果不存在就插入数据
						id, err := o.Insert(&item)
						fmt.Println(id)
						fmt.Println(err)
					}
				}
				res["data"] = id
				res["code"] = 200
				res["msg"] = "修改成功"
			}

		}
	}
	// 返回json格式给前端界面
	c.Data["json"] = res
	c.ServeJSON()
}

// 刪除商品属性组合
func (c *ProductController) DelAttr() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 构造数据查询条件
	var params ProductSpecs
	data := c.Ctx.Input.RequestBody
	//json数据封装到params对象中
	err := json.Unmarshal(data, &params)
	if err != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
	} else {
		if _, err := o.Delete(&params); err == nil {
			res["code"] = 200
			res["msg"] = "删除成功"
		} else {
			res["code"] = 401
			res["msg"] = "删除失败"
		}
	}

	// 返回json格式给前端界面
	c.Data["json"] = res
	c.ServeJSON()
}

// 获取商品属性表数据
func (c *ProductController) GetProductSpecs() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	paramsId := c.Ctx.Input.Param(":id")
	productId, _ := strconv.Atoi(paramsId)

	var ProductSpecsList []ProductSpecs
	num, err := o.QueryTable("ProductSpecs").Filter("product_id", productId).All(&ProductSpecsList)
	if err != nil {
		res["code"] = 401
	} else {
		res["total"] = num
	}
	res["code"] = 200
	res["msg"] = err
	res["data"] = ProductSpecsList
	c.Data["json"] = res
	c.ServeJSON()
}

// 刪除商品
func (c *ProductController) DelProduct() {
	res := make(map[string]interface{})
	o := orm.NewOrm()

	// 构造数据查询条件
	var params Product
	data := c.Ctx.Input.RequestBody
	//json数据封装到params对象中
	err := json.Unmarshal(data, &params)
	if err != nil {
		res["data"] = ""
		res["code"] = 401
		res["msg"] = err
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 更新商品IsDel字段
	num, sqlErr := o.Update(&params, "IsDel")

	if sqlErr != nil {
		res["data"] = ""
		res["code"] = 401
		res["msg"] = sqlErr
	} else {
		res["data"] = num
		res["code"] = 200
		res["msg"] = "删除成功"
	}
	c.Data["json"] = res
	c.ServeJSON()
}
