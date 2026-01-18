package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 在社区中创建帖子
// @Summary 在社区中创建帖子
// @Description 需要参数，通过结构体创建帖子信息
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.Post true "帖子json"
// @Security ApiKeyAuth
// @Success 1000 {object} controller.ResponseData
// @Router /post [post]
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

// GetPostDetailHandler 获取帖子详情
// @Summary 获取帖子详情
// @Description 需要参数，根据id查询帖子详情
// @Tags 帖子相关接口
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path string true "帖子id"
// @Security ApiKeyAuth
// @Success 1000 {object} models.ApiPostDetail
// @Router /post/:id [get]
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

// GetPostListHandler 获取所有帖子分页列表
// @Summary 获取所有帖子分页列表
// @Description 需要参数，根据分页参数查询帖子分页列表
// @Tags 帖子相关接口
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamPage true "帖子param"
// @Security ApiKeyAuth
// @Success 1000 {object} models.ApiPostPage
// @Router /posts [get]
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

// GetPostListHandler2 按时间排序 或 按分数排序列表 或综合热度排序，是否属于社区可选
// @Summary 按时间排序 或 按分数排序列表 或综合热度排序，是否属于社区可选
// @Description 需要参数，根据分页参数和排序参数和社区id(可选)查询帖子分页排序列表
// @Tags 帖子相关接口
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamPostListPage true "帖子param"
// @Security ApiKeyAuth
// @Success 1000 {object} models.ApiPostPage
// @Router /posts2 [get]
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
