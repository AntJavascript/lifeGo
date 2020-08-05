package models

import (
	"time"
)

// 注册表字段
type User struct {
	Id        int
	CreteTime time.Time
	Phone     string
	Password  string
	Email     string
	Face      string
	Name      string
}

// 请求的参数结构
type RegisterParams struct {
	Phone           string
	Password        string
	ConfirmPassword string
}
