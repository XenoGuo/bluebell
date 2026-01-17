package models

// 定义请求的参数结构体

// ParamSignUp 注册参数
type ParamSignUp struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RePasswd string `json:"re_passwd" binding:"required,eqfield=Password"`
}

// ParamLogin 登录参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamPage offset分页
type ParamPage struct {
	Page int64 `json:"page" form:"page" binding:"required"`
	Size int64 `json:"size" form:"size" binding:"required"`
}

// ParamPostListPage 通过发布时间/分数排序查询帖子分页列表
type ParamPostListPage struct {
	Page        int64  `json:"page" form:"page" binding:"required"`
	Size        int64  `json:"size" form:"size" binding:"required"`
	OrderBy     string `json:"order" form:"order" binding:"required,oneof='time' 'score' 'hot'"`
	CommunityID int64  `json:"community_id" form:"community_id"`
}

// ParamVote 投票参数
type ParamVote struct {
	PostID    int64 `json:"post_id,string" binding:"required"`        // 帖子ID
	Direction int8  `json:"direction,string" binding:"oneof= -1 0 1"` // 赞成票（1）返回票（-1）取消（0）
}
