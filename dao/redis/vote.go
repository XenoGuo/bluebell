package redis

import (
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneWeekInSeconds = 7 * 24 * 60 * 60
	scorePerVote     = 432  //每一票的分值
	decaySeconds     = 3600 //每3600秒（一小时）
)

var (
	ErrVoteTimeExpired = errors.New("投票时间已过")
	ErrVoteRepeat      = errors.New("不能投重复票")
)

/**
投票的几种情况：
direction=-1时，两种
	1. 之前没有投过票0，现在投反对票-1	差值绝对值1	-432
	2. 之前投赞成票1，现在改投反对票-1	2			-432*2
direction=0时，两种
	1. 之前投赞成票1，现在取消投票0		1			-432
	2. 之前投反对票-1，现在改投取消0	1			+ 432
direction=1时，有两种情况
	1. 之前没有投过票0，现在投赞成票1	1			+432
	2. 之前投反对票-1，现在改投赞成票1	2			+432*2

投票限制：
	每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票
		1. 到期之后将redis中保存的赞成票及反对票存储到mysql中
		2. 到期之后删除那个KeyPostVotedZSetPF
*/

func VoteForPost(userID string, postID string, nv float64) error {
	// 1. 判断投票的限制
	// 取redis取帖子发布时间
	postTime := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	// 超过一周不允许投票
	if (float64(time.Now().Unix()) - postTime) > oneWeekInSeconds {
		return ErrVoteTimeExpired
	}
	// 2. 更新分数
	// 当前用户之前的投票记录
	ov := rdb.ZScore(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Val()
	var op float64
	if ov == nv {
		return ErrVoteRepeat
	}
	if nv > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - nv) // 两次投票的绝对差值
	// 更新分数
	pipeline := rdb.TxPipeline()
	voteScore := op * diff * scorePerVote
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), voteScore, postID)
	hotScore := voteScore - (float64(time.Now().Unix())-postTime)/decaySeconds
	pipeline.ZIncrBy(getRedisKey(KeyPostHotZSet), hotScore, postID)
	// 3. 记录用户为该帖子投票的数据
	if nv == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postID), redis.Z{
			Score:  nv, // 当前是赞成或反对票
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
