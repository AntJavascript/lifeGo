package models

import (
	"fmt"
	"strconv"

	"encoding/json"
	"widget/tools"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 商品分类字段
type Classify struct {
	Id           int
	ClassifyName string
	ParentId     int
	Img          string
}
type JsonClass struct {
	Classify
	Children interface{}
}

// 查找二级分类函数
func findSubClass(Classifys Classify, id int) JsonClass {
	o := orm.NewOrm()
	var classify []Classify
	var jsonClass JsonClass

	o.QueryTable("classify").Filter("parent_id", id).All(&classify) // 获取一级商品分类
	jsonClass.Classify = Classifys
	jsonClass.Children = classify

	return jsonClass
}

type ClassifyListController struct {
	beego.Controller
}

// 查询分类
func (c *ClassifyListController) Get() {
	fmt.Println("LoginSession", c.GetSession("LoginSession"))
	o := orm.NewOrm()
	var classify []Classify
	var jsonClass = make([]interface{}, 0)

	o.QueryTable("classify").Filter("parent_id", 0).All(&classify) // 获取一级商品分类
	// 循环一级分类，查询二级分类
	for _, value := range classify {
		jsonClass = append(jsonClass, findSubClass(value, value.Id))
	}
	res := make(map[string]interface{})
	res["data"] = &jsonClass
	res["code"] = 200
	res["msg"] = ""
	res["total"] = len(jsonClass)
	c.Data["json"] = res
	c.ServeJSON()
}

// 新增分类
func (c *ClassifyListController) AddPost() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 获取ajax提交的参数
	ClassifyName := c.GetString("ClassifyName")
	ParentId := c.GetString("ParentId")
	ParentIdIntid, iderr := strconv.Atoi(ParentId)
	// 父级分类id错误
	if iderr != nil {
		c.Data["json"] = res
		res["msg"] = iderr
		c.ServeJSON()
		return
	}
	var filePath string
	// 处理分类图片
	f, h, err := c.GetFile("file")
	var parames Classify
	// 如果图片信息存在
	if err == nil {
		var uploadObj = make(map[string]string)
		uploadObj = tools.Upload(f, h)
		ok := c.SaveToFile("file", uploadObj["path"]) // 保存位置在 static/upload, 没有文件夹要先创建
		if ok != nil {
			res["code"] = 401
			res["msg"] = ok
			c.Data["json"] = res
			c.ServeJSON()
			return
		}
		parames.Img = uploadObj["path"]
	}

	parames.ParentId = ParentIdIntid
	parames.ClassifyName = ClassifyName

	id, err := o.Insert(&parames)

	if err != nil {
		res["data"] = ""
		res["code"] = 401
		res["msg"] = err
	} else {
		res["data"] = id
		res["code"] = 200
		res["msg"] = "新增成功"
		res["img"] = filePath
	}
	c.Data["json"] = res
	c.ServeJSON()
}

// 修改分类
func (c *ClassifyListController) EditPost() {
	var uploadObj = make(map[string]string)
	res := make(map[string]interface{})
	o := orm.NewOrm()
	var params Classify
	// 获取ajax提交的参数
	Id := c.GetString("Id")
	ClassifyName := c.GetString("ClassifyName")
	isId, err := strconv.Atoi(Id)

	// 处理分类图片
	f, h, err := c.GetFile("file")
	// 如果有图片信息
	if err == nil {
		uploadObj = tools.Upload(f, h)
		ok := c.SaveToFile("file", uploadObj["path"]) // 保存位置在 static/upload, 没有文件夹要先创建
		if ok != nil {
			res["code"] = 402
			res["msg"] = "上传失败"
			return
		}
		params.Img = uploadObj["path"]
	}
	params.Id = isId
	params.ClassifyName = ClassifyName
	if uploadObj["path"] != "" {
		_, err := o.Update(&params, "ClassifyName", "Img")
		if err != nil {
			fmt.Println(err)
			res["code"] = 401
			res["msg"] = err
		} else {
			res["code"] = 200
			res["msg"] = "修改成功"
			res["img"] = uploadObj["path"]
		}
	} else {
		_, err := o.Update(&params, "ClassifyName")
		if err != nil {
			fmt.Println(err)
			res["code"] = 401
			res["msg"] = err
		} else {
			res["code"] = 200
			res["msg"] = "修改成功"
			res["img"] = uploadObj["path"]
		}
	}

	c.Data["json"] = res
	c.ServeJSON()
}

// 根据id获取数据
func (c *ClassifyListController) ClassifyIdFindListData() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 获取ajax提交的参数
	paramsId := c.Ctx.Input.Param(":id")
	parseId, _ := strconv.Atoi(paramsId)
	//	params := Classify{Id: parseId}

	var List []Product
	num, err := o.QueryTable("Product").Filter("category_sub_id", parseId).Filter("ISDel", 0).Limit(10).All(&List)
	if err != nil {
		res["code"] = 401
		res["msg"] = err
	} else {
		res["code"] = 200
		res["msg"] = ""
		res["total"] = num
		res["data"] = List
	}
	// 返回json格式给前端界面
	c.Data["json"] = res
	c.ServeJSON()
}

// 刪除商品
func (c *ClassifyListController) DelClassify() {
	res := make(map[string]interface{})
	o := orm.NewOrm()

	// 构造数据查询条件
	var params Classify
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
	num, sqlErr := o.QueryTable("Classify").Filter("id", params.Id).Delete()

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
