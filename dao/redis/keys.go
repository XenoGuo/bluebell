package redis

// redis key注意使用命名空间的方式，方便业务
const (
	KeyPrefix = "bluebell:"

	// 全局榜单（ZSet）
	KeyPostTimeZSet  = "post:time"  // zset;帖子及发贴时间
	KeyPostScoreZSet = "post:score" // zset;帖子及投票的分数
	KeyPostHotZSet   = "post:hot"   // zset: 帖子的热度综合分数
	// key: community:{cid}:posts Set member=postID
	KeyPostVotedZSetPrefix = "post:voted:" // zset;用户投票的类型,参数是post_id

	// 社区帖子集合
	KeyCommunitySetPF = "community:" // set;保存每个分区下的帖子 (key:community:{cid}:posts)
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
