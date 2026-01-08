package controller

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//-------社区相关

func CommunityHandler(c *gin.Context) {
	// 查询所有的社区 (community_id,community_name) 以列表返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务端报错返回给前端
		return
	}
	ResponseSuccess(c, data)
}

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
