package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"path"
	"strconv"
	"time"

	"widget/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type CommonController struct {
	beego.Controller
}

//生成32位md5字串
//func GetMd5String(s string) string {
//	h := md5.New()
//	h.Write([]byte(s))
//	return hex.EncodeToString(h.Sum(nil))
//}

////生成Guid字串
//func UniqueId() string {
//	b := make([]byte, 48)

//	if _, err := io.ReadFull(rand.Reader, b); err != nil {
//		return ""
//	}
//	return GetMd5String(base64.URLEncoding.EncodeToString(b))
//}

// 文件上传
func (c *CommonController) Upload() {
	t := time.Now() // 当前时间
	res := make(map[string]interface{})
	f, h, err := c.GetFile("file")
	if f == nil || h == nil || err != nil {
		res["code"] = "401"
		res["msg"] = "没有上传文件"
		res["path"] = ""
		c.Data["json"] = res
		c.ServeJSON()
	}
	defer f.Close()
	var time = strconv.FormatInt(t.Unix(), 10) // 获取时间戳
	ext := path.Ext(h.Filename)                // 获取后缀名
	fmt.Println("上传图片地址：" + h.Filename)
	var filePath = "static/upload/" + time + ext // 文件路径名称
	ok := c.SaveToFile("file", filePath)         // 保存位置在 static/upload, 没有文件夹要先创建
	if ok != nil {
		res["code"] = 402
		res["msg"] = "上传失败"
		c.Data["json"] = res
		c.ServeJSON()
	}
	dbHost := beego.AppConfig.String("uploadPath")
	res["code"] = "200"
	res["path"] = dbHost + "/" + filePath
	c.Data["json"] = res
	c.ServeJSON()
}

// 评论文件上传
func (c *CommonController) CommentUpload() {
	res := make(map[string]interface{})
	loginSession := c.GetSession("LoginSession")
	if loginSession == nil {
		res["code"] = 401
		res["msg"] = "请先登录"
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	t := time.Now() // 当前时间

	f, h, err := c.GetFile("file")
	if f == nil || h == nil || err != nil {
		res["code"] = "401"
		res["msg"] = "没有上传文件"
		res["path"] = ""
		c.Data["json"] = res
		c.ServeJSON()
	}
	defer f.Close()
	var time = strconv.FormatInt(t.Unix()+rand.Int63(), 10) // 获取时间戳
	ext := path.Ext(h.Filename)                             // 获取后缀名
	fmt.Println("上传图片地址：" + h.Filename)
	var filePath = "static/comment/" + time + ext // 文件路径名称
	ok := c.SaveToFile("file", filePath)          // 保存位置在 static/comment, 没有文件夹要先创建
	if ok != nil {
		res["code"] = 402
		res["msg"] = "上传失败"
		c.Data["json"] = res
		c.ServeJSON()
	}
	dbHost := beego.AppConfig.String("uploadPath")
	res["code"] = "200"
	res["path"] = dbHost + "/" + filePath
	c.Data["json"] = res
	c.ServeJSON()
}

// 获取验证码
func (c *CommonController) GetverifyCode() {
	res := make(map[string]interface{})
	o := orm.NewOrm()
	var params models.SendMessage
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		res["msg"] = "参数解析错误"
		res["code"] = "400"
		res["data"] = ""
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if params.Phone == "" {
		res["msg"] = "手机号不能为空"
		res["code"] = "400"
		res["data"] = ""
	} else {
		var num int
		for {
			num = rand.Intn(10000)
			if num >= 1000 {
				break
			}
		}
		params.Code = num
		_, err := o.Insert(&params)
		if err != nil {
			res["msg"] = "插入失败"
			res["code"] = "401"
			res["data"] = "" // 生成4位数随机数
		} else {
			res["msg"] = ""
			res["code"] = "200"
			res["data"] = params.Code // 生成4位数随机数
		}

	}
	c.Data["json"] = res
	c.ServeJSON()
}
