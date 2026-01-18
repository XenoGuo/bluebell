package controller

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityHandler 查询所有的社区 (community_id,community_name) 以列表返回
// @Summary 查询所有的社区 (community_id,community_name) 以列表返回
// @Description 不需要参数，查询所有社区 的id 和 name 以列表json返回
// @Tags 社区相关接口
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 1000 {object} []models.Community
// @Router /community [get]
func CommunityHandler(c *gin.Context) {
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务端报错返回给前端
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 查询社区详情信息
// @Summary 查询社区详情信息
// @Description 需要参数，通过id查询社区详情
// @Tags 社区相关接口
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path string true "查询参数"
// @Security ApiKeyAuth
// @Success 1000 {object} models.CommunityDetail
// @Router /community/:id [get]
func CommunityDetailHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Warn("InvalidParam failed", zap.String("id", idStr), zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	if id <= 0 {
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 查询所有的社区 (community_id,community_name) 以列表返回
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务端报错返回给前端
		return
	}
	ResponseSuccess(c, data)
}
