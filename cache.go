package proxyinabox

import (
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var ipDomainCache *cache.Cache

type ipDomainCacheStruct struct {
	lock *sync.Mutex
	list map[string]struct{}
}

func init() {
	ipDomainCache = cache.New(30*time.Minute, 40*time.Minute)
}

//CheckIPDomain check domain num by ip
func CheckIPDomain(ip, domain string) bool {
	domains, has := ipDomainCache.Get(ip)
	if has {
		domains := domains.(ipDomainCacheStruct)
		domains.lock.Lock()
		defer domains.lock.Unlock()
		_, had := domains.list[domain]
		if had {
			return true
		}
		if len(domains.list) < DomainsPerIPHalfHour {
			domains.list[domain] = struct{}{}
			return true
		}
		return false
	}
	domains = ipDomainCacheStruct{
		list: make(map[string]struct{}),
		lock: new(sync.Mutex),
	}
	ipDomainCache.Set(ip, domains, cache.DefaultExpiration)
	return true
}
