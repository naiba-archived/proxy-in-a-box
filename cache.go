package proxyinabox

import (
	"sync"
	"sync/atomic"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var cacheInstance *cache.Cache

type domainList struct {
	lock *sync.Mutex
	list map[string]struct{}
}

func init() {
	cacheInstance = cache.New(30*time.Minute, 40*time.Minute)
}

//CheckIPLimit check ip limit
func CheckIPLimit(ip string) bool {
	num, ok := cacheInstance.Get(ip + "l")
	if ok {
		if *num.(*int32) > RequestLimitPerIPOneMinute {
			return false
		} else {
			atomic.AddInt32(num.(*int32), 1)
		}
	} else {
		var tmp int32
		tmp = 1
		cacheInstance.Set(ip+"l", &tmp, time.Minute*1)
	}
	return true
}

//CheckIPDomain check domain num by ip
func CheckIPDomain(ip, domain string) bool {
	domains, has := cacheInstance.Get(ip)
	if has {
		domains := domains.(domainList)
		domains.lock.Lock()
		defer domains.lock.Unlock()
		_, had := domains.list[domain]
		if had {
			return true
		}
		if len(domains.list) < DomainsPerIPHalfAnHour {
			domains.list[domain] = struct{}{}
			return true
		}
		return false
	}
	domains = domainList{
		list: make(map[string]struct{}),
		lock: new(sync.Mutex),
	}
	cacheInstance.Set(ip, domains, cache.DefaultExpiration)
	return true
}
