package controllers

import (
	"widget/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type MainController struct {
	beego.Controller
}

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(
		new(models.AdminUser),
		new(models.Classify),
		new(models.AttributeKey),
		new(models.AttributeValue),
		new(models.Product),
		new(models.Logistics),
		new(models.ProductSpecs),
		new(models.Address),
		new(models.SendMessage),
		new(models.User),
		new(models.Login),
		new(models.About),
		new(models.Article),
		new(models.Coupon),
		new(models.Cart),
		new(models.UserAddress),
		new(models.UserCoupon),
		new(models.UserOrder),
		new(models.UserComment),
	)
}
