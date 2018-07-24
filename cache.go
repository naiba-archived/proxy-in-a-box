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
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := Cache.Ping().Result()
	com.PanicIfNotNil(err)
}
