package controller

import (
	"bluebell/logic"
	"bluebell/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// PostVoteHandler 为指定的帖子投票
// @Summary 为指定的帖子投票
// @Description 需要参数，根据帖子id和投票类型进行投票
// @Tags 投票相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamVote true "帖子json"
// @Security ApiKeyAuth
// @Success 1000 {object} controller.ResponseData
// @Router /vote [post]
func PostVoteHandler(c *gin.Context) {
	// 参数处理
	p := new(models.ParamVote)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
		}
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}

	// 当前谁投票
	userID, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	err = logic.PostVote(userID, p)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
