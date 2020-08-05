package routers

import (
	"widget/controllers"
	"widget/models"

	"github.com/astaxie/beego"
)

func init() {
	// 首页
	beego.Router("/", &controllers.MainController{})
	// 登录
	beego.Router("/login", &models.LoginController{})
	// 获取用户列表
	beego.Router("/user/list", &models.WebUserController{}, "get:GetList")
	// 修改用户信息
	beego.Router("/user/EditUserInfo", &models.WebUserController{}, "post:EditUserInfo")
	// 获取用户信息
	beego.Router("/user/getUserInfo", &models.WebUserController{}, "get:GetUserInfo")
	// 上传文件
	beego.Router("/common/upload", &controllers.CommonController{}, "post:Upload")
	// 获取商品分类
	beego.Router("/classify/list", &models.ClassifyListController{})
	// 新增商品分类
	beego.Router("/classify/add", &models.ClassifyListController{}, "post:AddPost")
	// 删除商品分类
	beego.Router("/classify/del", &models.ClassifyListController{}, "post:DelClassify")
	// 新增商品分类Id获取商品
	beego.Router("/classifyIdFindListData/:id", &models.ClassifyListController{}, "get:ClassifyIdFindListData")
	// 修改商品分类
	beego.Router("/classify/edit", &models.ClassifyListController{}, "post:EditPost")
	// 新增商品属性
	beego.Router("/attributeKey/add", &models.AttributeKeyController{}, "post:AddPost")
	// 获取商品属性
	beego.Router("/attributeKey/list", &models.AttributeKeyController{}, "get:GetList")
	// 新增商品属性值
	beego.Router("/attributeValue/add", &models.AttributeValueController{}, "post:AddPost")
	// 获取商品属性值
	beego.Router("/attributeValue/list", &models.AttributeValueController{}, "get:GetList")
	// 获取商品
	beego.Router("/product/list", &models.ProductController{}, "get:GetList")
	// 根据id获取商品
	beego.Router("/product/detail/:id", &models.ProductController{}, "get:FindIdData")
	// 根据商品id获取属性
	beego.Router("/productSpecs/:id", &models.ProductController{}, "get:GetProductSpecs")
	// 新增商品
	beego.Router("/product/add", &models.ProductController{}, "post:AddPost")
	// 删除商品
	beego.Router("/product/del", &models.ProductController{}, "post:DelProduct")
	// 修改商品
	beego.Router("/product/edit/:id", &models.ProductController{}, "post:EditProductDetail")
	// 删除商品属性
	beego.Router("/product/delAttr", &models.ProductController{}, "post:DelAttr")
	// 获取物流列表
	beego.Router("/logistics/list", &models.LogisticsController{}, "get:GetList")
	// 新增物流
	beego.Router("/logistics/add", &models.LogisticsController{}, "post:AddPost")
	// 修改物流
	beego.Router("/logistics/edit", &models.LogisticsController{}, "post:EditPost")
	// 获取地址列表
	beego.Router("/address/list", &models.AddressController{})
	// 新增地址数据
	beego.Router("/address/add", &models.AddressController{}, "post:AddPost")
	// 修改地址数据
	beego.Router("/address/edit", &models.AddressController{}, "post:EditPost")
	// 生成随机数
	beego.Router("/common/GetverifyCode", &controllers.CommonController{}, "post:GetverifyCode")
	// 验证验证码的正确性
	beego.Router("/VerifyCode", &models.SendMessageController{}, "post:VerifyCode")
	// 注册
	beego.Router("/Register", &models.RegisterController{}, "post:Register")
	// 登录
	beego.Router("/AppLogin", &models.AppLoginController{}, "post:Login")
	// 退出
	beego.Router("/LoginOut", &models.AppLoginController{}, "post:LoginOut")
	// 修改密码
	beego.Router("/EditPassword", &models.EditpasswordController{}, "post:EditPassword")
	// 关于我们
	beego.Router("/About", &models.AboutController{}, "get:GetAuout")
	// 修改关于我们
	beego.Router("/EditAbout", &models.AboutController{}, "post:EditAbout")
	// 文章列表
	beego.Router("/Article/list", &models.ArticleController{}, "get:GetList")
	// 添加文章
	beego.Router("/Article/add", &models.ArticleController{}, "post:AddPost")
	// 文章详情
	beego.Router("/Article/detail/:id", &models.ArticleController{}, "get:GetDetail")
	// 修改文章
	beego.Router("/Article/edit/:id", &models.ArticleController{}, "post:EditArticle")
	// 优惠券列表
	beego.Router("/Coupon/list", &models.CouponController{}, "get:GetList")
	// 添加优惠券
	beego.Router("/Coupon/add", &models.CouponController{}, "post:AddPost")
	// 获取购物车列表
	beego.Router("/Cart/list", &models.CartController{}, "get:FindList")
	// 获取选中购物车列表
	beego.Router("/Cart/selectList", &models.CartController{}, "get:FindSelectList")
	// 添加购物车
	beego.Router("/Cart/add", &models.CartController{}, "post:AddCart")
	// 修改购物车数量
	beego.Router("/Cart/EditCartNumber", &models.CartController{}, "post:EditCartNumber")
	// 删除购物车
	beego.Router("/Cart/DelCart", &models.CartController{}, "post:DelCart")
	// 修改购物车状态
	beego.Router("/CartStatus", &models.CartController{}, "post:CartStatus")
	// 获取个人收货地址
	beego.Router("/User/GetUserAddressList", &models.UserAddressController{}, "post:GetUserAddressList")
	// 添加个人收货地址
	beego.Router("/User/AddUserAddress", &models.UserAddressController{}, "post:AddUserAddress")
	// 删除个人收货地址
	beego.Router("/User/DeleteUserAddress", &models.UserAddressController{}, "post:DeleteUserAddress")
	// 个人收货地址详情
	beego.Router("/User/UserAddressDetail/:id", &models.UserAddressController{}, "get:UserAddressDetail")
	// 修改个人收货地址
	beego.Router("/User/EditUserAddress", &models.UserAddressController{}, "post:EditUserAddress")
	// 获取个人优惠券
	beego.Router("/User/UserCouponList", &models.UserCouponController{}, "post:GetList")
	// 领取优惠券
	beego.Router("/User/ReceiveCoupon", &models.UserCouponController{}, "post:ReceiveCoupon")
	// 生成订单
	beego.Router("/Order/Buy", &models.OrderController{}, "post:Buy")
	// 后台展示订单列表
	beego.Router("/Order/List", &models.OrderController{}, "get:GetList")
	// 个人订单列表
	beego.Router("/User/Order/List", &models.OrderController{}, "get:GetUserList")
	// 获取订单详情
	beego.Router("/User/Order/Detail/:id", &models.OrderController{}, "get:GetUserOrderDetail")
	// 获取订单详情
	beego.Router("/Admin/Order/Detail/:id", &models.OrderController{}, "get:AdminGetUserOrderDetail")
	// 支付订单
	beego.Router("/User/Order/Pay", &models.OrderController{}, "post:Pay")
	// 取消订单
	beego.Router("/User/Order/CancelOrder", &models.OrderController{}, "post:CancelOrder")
	// 发货订单
	beego.Router("/Admin/Order/SendOrder", &models.OrderController{}, "post:SendOrder")
	// 确认收货
	beego.Router("/User/Order/Receiving", &models.OrderController{}, "post:Receiving")
	// 评论列表
	beego.Router("/User/Comment/GetList", &models.UserCommentController{}, "get:GetList")
	// 发布评论
	beego.Router("/User/Comment/Submit", &models.UserCommentController{}, "post:SubmitComment")
	// 上传评论图片
	beego.Router("/comment/upload", &controllers.CommonController{}, "post:CommentUpload")
	// 根据产品id获取评论列表
	beego.Router("/comment/list/:id", &models.UserCommentController{}, "get:FindProductIdGetComment")
}
