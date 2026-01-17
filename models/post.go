package models

import "time"

type Post struct {
	ID          int64     `db:"post_id" json:"id,string"`
	AuthorID    int64     `db:"author_id" json:"author_id,string"`
	CommunityID int64     `db:"community_id" json:"community_id,string" binding:"required"`
	Status      int32     `db:"status" json:"status"`
	Title       string    `db:"title" json:"title" binding:"required"`
	Content     string    `db:"content" json:"content" binding:"required"`
	CreateTime  time.Time `db:"create_time" json:"create_time"`
}

// ApiPostDetail 帖子详情接口的结构体
type ApiPostDetail struct {
	AuthorName string `json:"author_name"`
	VoteNum    int64  `json:"vote_num"`
	*Post
	*CommunityDetail `json:"community"`
}

type ApiPostPage struct {
	List  interface{} `json:"list"`
	Page  int64       `json:"page"`
	Size  int64       `json:"size"`
	Total int64       `json:"total"`
}
