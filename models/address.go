package models

import (
	"fmt"
	"strconv"

	"widget/tools"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 地址管理表字段
type Address struct {
	Id          int
	AddressName string
	ParentId    int
}
type ChildAddress struct {
	Address
	Children interface{}
}

// 查找二级分类函数
func findSubAddress(AddressArg Address, id int, isArea bool) ChildAddress {
	o := orm.NewOrm()
	var AddressList []Address
	var ChildAddress ChildAddress
	var ChildAddressArea = make([]interface{}, 0)

	o.QueryTable("Address").Filter("parent_id", id).All(&AddressList) // 获取上级数据
	fmt.Println(id, AddressList, isArea)
	ChildAddress.Address = AddressArg
	if !isArea {
		for _, value := range AddressList {
			ChildAddressArea = append(ChildAddressArea, findSubAddress(value, value.Id, true))
			fmt.Println(ChildAddressArea)
		}
		ChildAddress.Children = ChildAddressArea
	} else {
		ChildAddress.Children = AddressList
	}
	return ChildAddress
}

type AddressController struct {
	beego.Controller
}

// 查询分类
func (c *AddressController) Get() {
	o := orm.NewOrm()
	var AddressList []Address
	var ChildAddress = make([]interface{}, 0)

	o.QueryTable("Address").Filter("parent_id", 0).All(&AddressList) // 获取一级商品分类
	// 循环一级分类，查询二级分类
	for _, value := range AddressList {
		ChildAddress = append(ChildAddress, findSubAddress(value, value.Id, false))
	}
	res := make(map[string]interface{})
	res["data"] = &ChildAddress
	res["code"] = 200
	res["msg"] = ""
	res["total"] = len(ChildAddress)
	c.Data["json"] = res
	c.ServeJSON()
}

// 新增分类
func (c *AddressController) AddPost() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	// 获取ajax提交的参数
	AddressName := c.GetString("addressName")
	fmt.Println(AddressName)
	if AddressName == "" {
		res["msg"] = "参数不能为空"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	ParentId := c.GetString("parentId")
	ParentIdIntid, iderr := strconv.Atoi(ParentId)
	fmt.Println(ParentIdIntid)
	// 父级分类id错误
	if iderr != nil {
		res["msg"] = iderr
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	var params Address

	params.ParentId = ParentIdIntid
	params.AddressName = AddressName

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
	c.Data["json"] = res
	c.ServeJSON()
}
func (c *AddressController) EditPost() {
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
