package redis

// redis key注意使用命名空间的方式，方便业务
const (
	KeyPrefix              = "bluebell:"
	KeyPostTimeZSet        = "post:time"   // zset;帖子及发贴时间
	KeyPostScoreZSet       = "post:score"  // zset;帖子及投票的分数
	KeyPostVotedZSetPrefix = "post:voted:" // zset;用户投票的类型,参数是post_id
)
