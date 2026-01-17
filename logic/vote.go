package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"errors"
	"strconv"

	"go.uber.org/zap"
)

var (
	ErrVoteWarn = errors.New("投票重复")
	ErrVoteErr  = errors.New("只允许在投票期间进行投票")
)

// 基于用户投票相关算法：https://www.ruanyifeng.com/blog/algorithm/

// 简化版的投票算法
// 投一票就加432分 86400/200 -> 200张赞成票可以给帖子续一天 -> 《redis实战》
/**
投票的几种情况：
direction=-1时，两种
	1. 之前没有投过票，现在投反对票
	2. 之前投赞成票，现在改投反对票
direction=0时，两种
	1. 之前没有投过票，现在取消投票
	2. 之前投反对票，现在改投取消
direction=1时，有两种情况
	1. 之前没有投过票，现在投赞成票
	2. 之前投反对票，现在改投赞成票

投票限制：
	每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票
		1. 到期之后将redis中保存的赞成票及反对票存储到mysql中
		2. 到期之后删除那个KeyPostVotedZSetPF
*/

// PostVote 为帖子投票
func PostVote(userID int64, p *models.ParamVote) error {
	//判断投票的限制
	err := redis.VoteForPost(strconv.FormatInt(userID, 10), strconv.FormatInt(p.PostID, 10), float64(p.Direction))
	if errors.Is(err, redis.ErrVoteTimeExpired) {
		zap.L().Error("PostVote err", zap.Error(err))
		zap.L().Warn("不允许在投票期间之外投票")
		return ErrVoteErr
	}
	if errors.Is(err, redis.ErrVoteRepeat) {
		zap.L().Error("PostVote err", zap.Error(err))
		zap.L().Warn("不能重复投相同票")
		return ErrVoteWarn
	}
	return nil
}
