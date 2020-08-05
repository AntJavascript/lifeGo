package tools

import (
	"math/rand"
	"mime/multipart"
	"path"
	"strconv"
	"time"
)

// 上传图片
func Upload(f multipart.File, h *multipart.FileHeader) map[string]string {
	t := time.Now()
	res := make(map[string]string)
	if f == nil || h == nil {
		res["code"] = "401"
		res["msg"] = "没有上传文件"
		res["path"] = ""
		return res
	}
	defer f.Close()
	var time = strconv.FormatInt(t.Unix(), 10)   // 获取时间戳
	ext := path.Ext(h.Filename)                  // 获取后缀名
	var filePath = "static/upload/" + time + ext // 文件路径名称
	res["code"] = "200"
	res["path"] = filePath
	return res
}

// 统一返回json格式处理
func ResultData(map[string]interface{}) {
}

// 上传评论图片
func CommentUpload(f multipart.File, h *multipart.FileHeader) map[string]string {
	t := time.Now()
	res := make(map[string]string)
	if f == nil || h == nil {
		res["code"] = "401"
		res["msg"] = "没有上传文件"
		res["path"] = ""
		return res
	}
	defer f.Close()
	var time = strconv.FormatInt(t.Unix()+rand.Int63(), 10) // 获取时间戳
	ext := path.Ext(h.Filename)                             // 获取后缀名
	var filePath = "static/comment/" + time + ext           // 文件路径名称
	res["code"] = "200"
	res["path"] = filePath
	return res
}
