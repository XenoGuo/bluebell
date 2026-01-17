package redis

import (
	"bluebell/models"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	INITSCORE  = 1
	ORDERTIME  = "time"
	ORDERSCORE = "score"
	ORDERHOT   = "hot"
)

// CreatePost 创建帖子
func CreatePost(postID int64, CommunityID int64) error {
	pipeline := rdb.TxPipeline()

	// 格式化参数
	postIDStr := strconv.FormatInt(postID, 10)
	now := float64(time.Now().Unix())
	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  now,
		Member: postIDStr,
	})
	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  INITSCORE,
		Member: postIDStr,
	})
	pipeline.ZAdd(getRedisKey(KeyPostHotZSet), redis.Z{
		Score:  INITSCORE,
		Member: postIDStr,
	})
	communitySetKey := getRedisKey(KeyCommunitySetPF) + strconv.FormatInt(CommunityID, 10)
	pipeline.SAdd(communitySetKey, postIDStr)
	_, err := pipeline.Exec()
	return err
}

// GetPostIDsInOrder 根据 orderBy 的条件查询postID切片
func GetPostIDsInOrder(p *models.ParamPostListPage) ([]string, error) {
	key := getRedisKey(KeyPostTimeZSet)
	// 根据用户请求中携带的排序方式 确定要查的redis key
	if p.OrderBy == ORDERSCORE {
		key = getRedisKey(KeyPostScoreZSet)
	} else if p.OrderBy == ORDERHOT {
		key = getRedisKey(KeyPostHotZSet)
	}
	// 确定查询起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	// ZREVRANGE 查询,
	return rdb.ZRevRange(key, start, end).Result()
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的票数
func GetPostVoteData(ids []string) (data []int64, err error) {
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	for _, cmder := range cmders {
		cmd := cmder.(*redis.IntCmd).Val()
		data = append(data, cmd)
	}
	return
}

func GetCommunityPostIDsOrder(p *models.ParamPostListPage) ([]string, error) {
	// 1. 选择排序用的全局ZSet
	var orderKey string
	switch p.OrderBy {
	case ORDERSCORE:
		orderKey = getRedisKey(KeyPostScoreZSet)
	case ORDERHOT:
		orderKey = getRedisKey(KeyPostHotZSet)
	default:
		orderKey = getRedisKey(KeyPostTimeZSet)
	}
	// 2. 社区过滤集合 community:xxxx
	communityKey := getRedisKey(fmt.Sprintf("%s%d", KeyCommunitySetPF, p.CommunityID))
	// 3. 临时交集key
	tempKey := getRedisKey(fmt.Sprintf("tmp:community:%d:%s", p.CommunityID, p.OrderBy))
	// 4. 计算分页范围
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	// 5. 交集+过期+分页查询（一个pipeline）
	// 缓存命中
	if rdb.Exists(tempKey).Val() > 0 {
		return rdb.ZRevRange(tempKey, start, end).Result()
	}

	// 缓存未命中
	pipeline := rdb.Pipeline()
	pipeline.ZInterStore(tempKey, redis.ZStore{
		Aggregate: "MAX",
	}, communityKey, orderKey)
	pipeline.Expire(tempKey, time.Minute)
	cmd := pipeline.ZRevRange(tempKey, start, end)
	_, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	return cmd.Result()
}
