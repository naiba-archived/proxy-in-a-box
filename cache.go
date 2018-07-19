package proxyinabox

import (
	"sync"
	"sync/atomic"
	"time"

	cache "github.com/patrickmn/go-cache"
)

type domainList struct {
	lock *sync.Mutex
	list map[string]struct{}
}

//Caches system cache
type Caches struct {
	cacheInstance *cache.Cache
	domainCache   sync.Map
	proxyCache    sync.Map
	proxyIndex    sync.Map
}

//TODO: init

//CacheInstance system cache instance
var CacheInstance *Caches

//CheckIPLimit check ip limit
func (c *Caches) CheckIPLimit(ip string) bool {
	num, ok := c.cacheInstance.Get(ip + "l")
	if ok {
		if *num.(*int32) > Config.Sys.RequestLimitPerIP {
			return false
		}
		atomic.AddInt32(num.(*int32), 1)
	} else {
		var tmp int32
		tmp = 1
		c.cacheInstance.Set(ip+"l", &tmp, time.Minute)
	}
	return true
}

//CheckIPDomain check domain num by ip
func (c *Caches) CheckIPDomain(ip, domain string) bool {
	domains, has := c.cacheInstance.Get(ip)
	if has {
		domains := domains.(domainList)
		domains.lock.Lock()
		defer domains.lock.Unlock()
		_, had := domains.list[domain]
		if had {
			return true
		}
		if len(domains.list) < Config.Sys.DomainsPerIP {
			domains.list[domain] = struct{}{}
			return true
		}
		return false
	}
	domains = domainList{
		list: make(map[string]struct{}),
		lock: new(sync.Mutex),
	}
	c.cacheInstance.Set(ip, domains, cache.DefaultExpiration)
	return true
}

//GetProxyByURI get proxy by uri string
func (c *Caches) GetProxyByURI(ps string) (Proxy, bool) {
	p, has := c.proxyIndex.Load(ps)
	return p.(Proxy), has
}

//SaveProxy save a proxy
func (c *Caches) SaveProxy(p Proxy) (e error) {
	if e = DB.Save(&p).Error; e != nil {
		return
	}
	c.proxyIndex.Store(p.URI(), p)
	c.proxyCache.Store(p.ID, p)
	return
}
