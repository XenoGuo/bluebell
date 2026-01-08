package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"errors"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) error {
	// 1.生成post id
	p.ID = snowflake.GenID()
	// 2.保存到数据库
	if err := mysql.CreatePost(p); err != nil {
		return err
	}
	// 3.返回
	return nil
}

// GetPostDetailById 根据帖子id返回帖子详情
func GetPostDetailById(id int64) (*models.ApiPostDetail, error) {
	post, err := mysql.GetPostById(id)
	if err != nil {
		if errors.Is(err, mysql.ErrNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		if errors.Is(err, mysql.ErrNotFound) {
			return nil, ErrUserNotExists
		}
		return nil, err
	}
	communityDetail, err := mysql.GetCommunityDetailById(post.CommunityID)
	if err != nil {
		return nil, err
	}
	postDetail := new(models.ApiPostDetail)
	postDetail.AuthorName = user.Username
	postDetail.CommunityDetail = communityDetail
	postDetail.Post = post
	return postDetail, nil
}

// GetPostList 返回帖子分页列表
func GetPostList(page, size int64) (data *models.ApiPostPage, err error) {
	total, err := mysql.GetPostTotalCount()
	if err != nil {
		return nil, err
	}
	offset := (page - 1) * size
	list, err := mysql.GetPostList(size, offset)
	if err != nil {
		return nil, err
	}
	d := &models.ApiPostPage{
		Total: total,
		List:  list,
		Page:  page,
		Size:  size,
	}
	return d, nil
}
