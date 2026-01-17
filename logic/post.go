package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"errors"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) error {
	// 1.生成post id
	p.ID = snowflake.GenID()
	// 2.保存到数据库
	if err := mysql.CreatePost(p); err != nil {
		return err
	}
	// 记录帖子信息到redis
	if err := redis.CreatePost(p.ID, p.CommunityID); err != nil {
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
	postDetail := &models.ApiPostDetail{
		AuthorName:      user.Username,
		CommunityDetail: communityDetail,
		Post:            post,
	}
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

func GetPostListSorted(p *models.ParamPostListPage) (*models.ApiPostPage, error) {
	if p.CommunityID > 0 {
		sort, err := getCommunityPostsBySort(p)
		if err != nil {
			return nil, err
		}
		return sort, nil
	} else {
		sort, err := getPostListBySort(p)
		if err != nil {
			return nil, err
		}
		return sort, nil
	}
}

func getPostListBySort(p *models.ParamPostListPage) (*models.ApiPostPage, error) {
	return getPostsByGetter(
		func() ([]string, error) { return redis.GetPostIDsInOrder(p) },
		p.Page, p.Size,
		"post list is empty",
	)
}

func getCommunityPostsBySort(p *models.ParamPostListPage) (*models.ApiPostPage, error) {
	return getPostsByGetter(
		func() ([]string, error) { return redis.GetCommunityPostIDsOrder(p) },
		p.Page, p.Size,
		"community post list is empty",
		zap.Int64("community_id", p.CommunityID),
	)
}

func buildPostPage(ids []string, page, size int64) (*models.ApiPostPage, error) {
	posts, err := mysql.GetPostsByIDs(ids)
	if err != nil {
		return nil, err
	}

	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}

	list := make([]models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			continue
		}
		c, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			continue
		}

		list = append(list, models.ApiPostDetail{
			AuthorName:      user.Username,
			CommunityDetail: c,
			Post:            post,
			VoteNum:         voteData[idx],
		})
	}

	return &models.ApiPostPage{
		List:  list,
		Page:  page,
		Size:  size,
		Total: int64(len(list)),
	}, nil
}

type idsGetter func() ([]string, error)

func getPostsByGetter(getIDs idsGetter, page, size int64, emptyLog string, fields ...zap.Field) (*models.ApiPostPage, error) {
	ids, err := getIDs()
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn(emptyLog, fields...)
		return nil, nil
	}
	return buildPostPage(ids, page, size)
}
