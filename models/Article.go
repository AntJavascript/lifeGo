package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Article struct {
	Id         int
	Title      string
	Desc       string
	Content    string
	CreateTime time.Time
	ProductId  int
	ViewCount  int
}

type ArticleController struct {
	beego.Controller
}

// 获取文章列表
func (c *ArticleController) GetList() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 获取limit ，每页需要展示的数量
	var pasize = c.GetString("limit")
	limit, err := strconv.ParseInt(pasize, 10, 64)
	// 当前处于第几页
	var page = c.GetString("page")
	currentPage, err := strconv.ParseInt(page, 10, 64)

	var list []Article
	var json = make([]interface{}, 0)
	num, err := o.QueryTable("Article").Limit(limit, currentPage-1).All(&list)
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

// 新增文章
func (c *ArticleController) AddPost() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 构造数据查询条件
	var params Article
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

// 文章详情
func (c *ArticleController) GetDetail() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 获取ajax提交的参数
	paramsId := c.Ctx.Input.Param(":id")
	parseId, _ := strconv.Atoi(paramsId)
	params := Article{Id: parseId}

	err := o.Read(&params)
	if err != nil {
		res["code"] = 401
		res["msg"] = err
	} else {
		res["code"] = 200
		res["msg"] = ""
		res["data"] = params
	}
	// 返回json格式给前端界面
	c.Data["json"] = res
	c.ServeJSON()
}

// 修改详情
func (c *ArticleController) EditArticle() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 构造数据查询条件
	var params Article
	data := c.Ctx.Input.RequestBody
	//json数据封装到params对象中
	err := json.Unmarshal(data, &params)
	if err != nil {
		res["data"] = ""
		res["code"] = 400
		res["msg"] = err
	} else {
		// 插入数据库, 返回id和err错误
		id, err := o.Update(&params)
		if err != nil {
			res["data"] = ""
			res["code"] = 401
			res["msg"] = err
		} else {
			res["data"] = id
			res["code"] = 200
			res["msg"] = "修改成功"
		}
	}
	// 返回json格式给前端界面
	c.Data["json"] = res
	c.ServeJSON()
}
