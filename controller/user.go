package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 注册用户
// @Summary 注册用户
// @Description 需要参数，根据用户名  密码和确认密码 注册用户
// @Tags 注册相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamSignUp true "帖子json"
// @Security
// @Success 1000 {object} controller.ResponseData
// @Router /signup [get]
func SignUpHandler(c *gin.Context) {
	// 1. 参数校验
	var p = new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("SignUpHandler ShouldBindJSON err", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2. 业务处理
	if err := logic.SignUp(p); err != nil {
		ResponseError(c, CodeUserExist)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// LoginHandler 用户登录
// @Summary 用户登录
// @Description 需要参数，根据用户名和密码登录
// @Tags 登录相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamLogin true "帖子json"
// @Security
// @Success 1000 {object} controller.ResponseData
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	// 获取参数和参数校验
	var p = new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("LoginHandler ShouldBindJSON err", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 业务逻辑处理
	token, err := logic.Login(p)
	if err != nil {
		zap.L().Error("LoginHandler err", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, logic.ErrUserNotExists) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}
	// 返回response
	ResponseSuccess(c, token)
}
