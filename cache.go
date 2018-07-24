package proxyinabox

import (
	"github.com/go-redis/redis"
	"github.com/naiba/com"
)

//Cache cache instance
var Cache *redis.Client

func initRedis() {
	Cache = redis.NewClient(&redis.Options{
		Addr:     Config.Redis.Host + ":" + Config.Redis.Port,
		Password: Config.Redis.Pass, // no password set
		DB:       Config.Redis.Db,   // use default DB
	})
	_, err := Cache.Ping().Result()
	com.PanicIfNotNil(err)
}
