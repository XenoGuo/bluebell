package models

// 定义请求的参数结构体

// 注册
type ParamSignUp struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RePasswd string `json:"re_passwd" binding:"required,eqfield=Password"`
}

// 登录
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// offset分页
type PageInfo struct {
	Page int64 `form:"page" binding:"required"`
	Size int64 `form:"size" binding:"required"`
}
