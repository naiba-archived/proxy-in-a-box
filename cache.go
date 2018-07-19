package proxyinabox

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	cache "github.com/patrickmn/go-cache"
)

const ipLimitPrefix = "il"

type liveDomain struct {
	lastActive int64
	proxies    sync.Map //map[uint]time.Unix
}

var cacheInstance *cache.Cache

var liveDomains sync.Map //map[domain]*liveDomain

var proxyCache sync.Map //map[id]*Proxy
var proxyIndex sync.Map //map[string]id
var proxyQueue struct {
	mu   sync.Mutex
	list []uint
}

//Caches system cache
type Caches struct{}

//CacheInstance cache instance
var CacheInstance Caches

//CheckIPLimit check ip limit
func (c Caches) CheckIPLimit(ip string) bool {
	num, ok := cacheInstance.Get(ipLimitPrefix + ip)
	if ok {
		if *num.(*int32) > Config.Sys.RequestLimitPerIP {
			return false
		}
		atomic.AddInt32(num.(*int32), 1)
	} else {
		var tmp int32
		tmp = 1
		cacheInstance.Set(ipLimitPrefix+ip, &tmp, time.Minute)
	}
	return true
}

//CheckIPDomain check domain num by ip
func (c Caches) CheckIPDomain(ip, domain string) bool {
	type domainList struct {
		lock *sync.Mutex
		list map[string]struct{}
	}
	domains, has := cacheInstance.Get(ip)
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
	cacheInstance.Set(ip, domains, cache.DefaultExpiration)
	return true
}

//GetProxyIDByURI get proxy by uri string
func (c Caches) GetProxyIDByURI(ps string) (interface{}, bool) {
	return proxyIndex.Load(ps)
}

//GetFreshProxy dispatch proxy
func (c Caches) GetFreshProxy(domain string) (*Proxy, error) {
	proxyQueue.mu.Lock()
	defer proxyQueue.mu.Unlock()
	if len(proxyQueue.list) == 0 {
		return nil, errors.New(fmt.Sprint("has no proxy in system."))
	}
	for i := 0; i < len(proxyQueue.list); i++ {
		// load proxy
		p := proxyQueue.list[i]
		ld, _ := liveDomains.LoadOrStore(domain, &liveDomain{lastActive: time.Now().Unix()})
		now := time.Now().Unix()
		if now-ld.(*liveDomain).lastActive < 4 {
			// domain is just active
			t, has := ld.(*liveDomain).proxies.Load(p)
			if has && now-t.(int64) < 3 {
				// proxy is juest used
				fmt.Println("proxy is juest used", p, now, t)
				continue
			}
		}
		ld.(*liveDomain).proxies.Store(p, now)
		ld.(*liveDomain).lastActive = now
		pp, has := proxyCache.Load(p)
		if !has {
			return nil, errors.New("lost proxy cache")
		}
		//swap used proxy to last
		if i > 0 {
			proxyQueue.list = append(proxyQueue.list[1:i], proxyQueue.list[i+1:]...)
			proxyQueue.list = append(proxyQueue.list, p)
		} else {
			proxyQueue.list = append(proxyQueue.list[1:], proxyQueue.list[0])
		}
		return pp.(*Proxy), nil
	}
	return nil, errors.New("has no free proxy")
}

//SaveProxy save a proxy
func (c Caches) SaveProxy(p Proxy) (e error) {
	if e = DB.Save(&p).Error; e != nil {
		return
	}
	proxyIndex.Store(p.URI(), p.ID)
	proxyCache.Store(p.ID, &p)
	proxyQueue.mu.Lock()
	defer proxyQueue.mu.Unlock()
	proxyQueue.list = append(proxyQueue.list, p.ID)
	return
}

//DeleteProxy save a proxy
func (c Caches) DeleteProxy(p Proxy) (e error) {
	if e = DB.Delete(&p).Error; e != nil {
		return
	}
	proxyIndex.Delete(p.URI())
	proxyCache.Delete(p.ID)
	proxyQueue.mu.Lock()
	defer proxyQueue.mu.Unlock()
	for i, pq := range proxyQueue.list {
		if pq == p.ID {
			proxyQueue.list = append(proxyQueue.list[:i], proxyQueue.list[i+1:]...)
			return
		}
	}
	fmt.Println("delete invalid proxy", p)
	return
}

//clear unused cache
func clearCacheWorker() {
	go func() {
		for {
			clearItem := 0
			now := time.Now().Unix()
			liveDomains.Range(func(key, value interface{}) bool {
				if now-value.(*liveDomain).lastActive > 3 {
					liveDomains.Delete(key)
					clearItem++
				} else {
					value.(*liveDomain).proxies.Range(func(k1, v1 interface{}) bool {
						if now-v1.(int64) > 3 {
							value.(*liveDomain).proxies.Delete(k1)
							clearItem++
						}
						return true
					})
				}
				return true
			})
			fmt.Println("clearCacheWorker", clearItem)
			time.Sleep(time.Minute)
		}
	}()
}
