package proxyinabox

import (
	"errors"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/naiba/com"
)

var cache *redis.Client

//TCache cache instance struct
type TCache struct{}

//CI cache instance
var CI TCache

func initRedis() {
	cache = redis.NewClient(&redis.Options{
		Addr:     Config.Redis.Host + ":" + Config.Redis.Port,
		Password: Config.Redis.Pass, // no password set
		DB:       Config.Redis.Db,   // use default DB
	})
	_, err := cache.Ping().Result()
	com.PanicIfNotNil(err)
}

//PickProxy get a fresh proxy
func (c TCache) PickProxy(host string) (string, error) {
	var p string
	var err error
	cache.Pipelined(func(pipe redis.Pipeliner) error {
		var ps []string
		ps, err = pipe.Sort("pp-*", &redis.Sort{Order: "ASC", Count: 10}).Result()
		if err != nil || len(ps) < 1 {
			err = errors.New("Proxy IP is not in stock")
			return nil
		}
		for i := 0; i < len(ps); i++ {
			ldkey := "ld-" + host + ps[i]
			if pipe.Get(ldkey).Err() == nil {
				pipe.Set(ldkey, nil, time.Second*3)
				ps = strings.Split(ps[i], "-")
				if len(ps) != 3 {
					err = errors.New("Unable to resolve proxy")
					return nil
				}
				p = ps[1] + ":" + ps[2]
				break
			}
		}
		return nil
	})
	if len(p) == 0 {
		err = errors.New("Proxy IP pool is too busy")
	}
	return p, err
}

//IPLimiter limit ip
func (c TCache) IPLimiter(ip string) bool {
	var key = "ipl-" + ip
	var count int64
	var err error
	cache.Pipelined(func(pipe redis.Pipeliner) error {
		count, err = pipe.Incr(key).Result()
		if err != nil {
			pipe.Set(key, 1, time.Second).Err()
			return nil
		}
		pipe.Expire(key, time.Second)
		return nil
	})
	return count <= Config.Sys.RequestLimitPerIP
}

//HostLimiter host limiter
func (c TCache) HostLimiter(ip, host string) bool {
	var key = "hl-" + ip
	count := 0
	cache.Pipelined(func(pipe redis.Pipeliner) error {
		now := time.Now().Unix()
		pipe.HSet(key, host, now)
		pipe.Expire(key, time.Minute*30)
		for _, h := range pipe.HKeys(key).Val() {
			last, _ := pipe.HGet(key, h).Int64()
			if now-last < 60*60*30 {
				count++
			} else {
				pipe.HDel(key, h)
			}
		}
		return nil
	})
	return count <= 10
}
