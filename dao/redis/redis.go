package redis

import (
	"bluebell/settings"
	"fmt"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

func Init(cfg *settings.RedisConfig) (err error) {
	//rdb = redis.NewClient(&redis.Options{
	//	Addr: fmt.Sprintf("%s:%d",
	//		viper.GetString("redis.host"),
	//		viper.GetInt("redis.port"),
	//	),
	//	Password: viper.GetString("redis.password"),
	//	DB:       viper.GetInt("redis.db"),
	//	PoolSize: viper.GetInt("redis.pool_size"),
	//})

	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	_ = rdb.Close()
}
