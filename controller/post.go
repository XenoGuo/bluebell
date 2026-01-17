package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreatePostHandler(c *gin.Context) {
	// 1.获取参数与参数校验
	// c.ShouldBindJSON() // validator --> binding tag
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreatePost ShouldBindJSON err", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	userId, err := getCurrentUser(c)
	if err != nil {
		zap.L().Error("CreatePost GetCurrentUser err", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userId
	// 2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost err", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, nil)
}

func GetPostDetailHandler(c *gin.Context) {
	// 1.获取参数 （从url中获取帖子的id）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("GetPostDetail err", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2.根据id取帖子数据
	data, err := logic.GetPostDetailById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostDetailById err", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, data)
}

func GetPostListHandler(c *gin.Context) {
	// 1.获取分页参数
	p := new(models.ParamPage)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostList ShouldBindQuery err", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Size <= 0 {
		p.Size = 10
	}
	if p.Size > 50 {
		p.Size = 50
	}
	// 2.根据参数取分页数据
	posts, err := logic.GetPostList(p.Page, p.Size)
	if err != nil {
		zap.L().Error("logic.GetPostListAll err", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.响应数据
	ResponseSuccess(c, posts)
}

// GetPostListHandler2 按时间排序 或 按分数排序列表 或综合热度排序
func GetPostListHandler2(c *gin.Context) {
	// 1.获取分页参数
	p := new(models.ParamPostListPage)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostList ShouldBindQuery err", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Size <= 0 {
		p.Size = 10
	}
	if p.Size > 50 {
		p.Size = 50
	}
	// 2.根据参数取分页数据
	posts, err := logic.GetPostListSorted(p)
	if err != nil {
		zap.L().Error("logic.GetPostListSorted err", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.响应数据
	ResponseSuccess(c, posts)
}
