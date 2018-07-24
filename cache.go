package proxyinabox

import (
	"errors"
	"fmt"
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
	var ps []string

	ps = cache.ZRange("ppc", 0, 10).Val()
	if len(ps) < 1 {
		return "", errors.New("Proxy IP is not in stock")
	}
	for i := 0; i < len(ps); i++ {
		//find in liveDomains
		ldkey := "ld" + host + ps[i]

		if cache.Get(ldkey).Err() == redis.Nil {
			//set used
			cache.Set(ldkey, nil, time.Second*3)
			//update proxy used time
			cache.ZIncr("ppc", redis.Z{Score: 1, Member: ps[1]})
			//get proxy
			p = cache.HGet("ppi", ps[i]).Val()
			break
		}
	}
	if len(p) == 0 {
		return "", errors.New("Proxy IP pool is too busy")
	}
	return p, nil
}

//IPLimiter limit ip
func (c TCache) IPLimiter(ip string) bool {
	var key = "ipl" + ip
	var count int64
	cache.Expire(key, time.Second)
	count = cache.Incr(key).Val()
	return count <= Config.Sys.RequestLimitPerIP
}

//HostLimiter host limiter
func (c TCache) HostLimiter(ip, host string) bool {
	var key = "hl" + ip
	count := 0
	now := time.Now().Unix()
	cache.HSet(key, host, now)
	cache.Expire(key, time.Minute*30)
	for _, h := range cache.HKeys(key).Val() {
		last, _ := cache.HGet(key, h).Int64()
		if now-last < 60*60*30 {
			count++
		} else {
			cache.HDel(key, h)
		}
	}
	return count <= Config.Sys.DomainsPerIP
}

//HasProxy has proxy
func (c TCache) HasProxy(p string) bool {
	return cache.HGet("ppr", p).Err() == nil
}

//SaveProxy save proxy
func (c TCache) SaveProxy(p Proxy) error {
	DB.Save(&p)
	cache.ZAdd("ppc", redis.Z{Member: p.ID})
	cache.HSet("ppi", string(p.ID), fmt.Sprintf("%s:%s", p.IP, p.Port))
	cache.HSet("ppr", fmt.Sprintf("%s:%s", p.IP, p.Port), p.ID)
	return nil
}

//DeleteProxy save proxy
func (c TCache) DeleteProxy(p Proxy) {
	if p.ID > 0 {
		cache.ZRem("ppc", redis.Z{Member: p.ID})
		cache.HDel("ppi", string(p.ID))
		cache.HDel("ppr", fmt.Sprintf("%s:%s", p.IP, p.Port))

		DB.Delete(&p)
	}
}
