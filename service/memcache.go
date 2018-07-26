package service

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/naiba/proxyinabox"
)

/*
==============
    代理池
==============
*/
type proxyEntry struct {
	p *proxyinabox.Proxy
	n int64
	l sync.Mutex
}

type sortableProxyList []*proxyEntry

func (p sortableProxyList) Len() int {
	return len(p)
}

func (p sortableProxyList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p sortableProxyList) Less(i, j int) bool {
	return p[i].n < p[j].n
}

type proxyList struct {
	l  sync.Mutex
	pl []*proxyEntry
}

/*
==============
   域名IP池
==============
*/
type domainScheduling struct {
	l  sync.Mutex
	dl map[string]*proxyList
}

/*
==============
   IP计数池
==============
*/
type ipActivityEntry struct {
	l          sync.Mutex
	lastActive int64
	num        int64
}
type ipActivity struct {
	l    sync.Mutex
	list map[string]*ipActivityEntry
}

//MemCache memory cache
type MemCache struct {
	proxies *proxyList
	domains *domainScheduling
	ips     *ipActivity
}

//TODO:new MemCache

//PickProxy rt
func (c *MemCache) PickProxy(req *http.Request) (string, error) {
	c.proxies.l.Lock()
	defer c.proxies.l.Unlock()

	length := len(c.proxies.pl)
	domain := req.Host
	now := time.Now().Unix()
	var candidate map[string]struct{}
	if length == 0 {
		return "", fmt.Errorf("%s", "There is no proxy in the proxy pool.")
	}
	sort.Sort(sortableProxyList(c.proxies.pl))
	candidate = make(map[string]struct{})
	c.domains.l.Lock()
	defer c.domains.l.Unlock()
	if pl, has := c.domains.dl[domain]; has {
		pl.l.Lock()
		defer pl.l.Unlock()
		sort.Sort(sortableProxyList(pl.pl))

		//清理长久未活动的代理
		for i, p := range pl.pl {
			p.l.Lock()
			defer p.l.Unlock()
			if p.n-now < 3 {
				candidate[p.p.IP] = struct{}{}
			} else {
				pl.pl = append(pl.pl[:i], pl.pl[i+1:]...)
			}
		}
	}

	for i := 0; i < length; i++ {
		if _, has := candidate[c.proxies.pl[i].p.IP]; !has {
			var schema string
			if c.proxies.pl[i].p.NotHTTPS {
				schema = "http://"
			} else {
				schema = "https://"
			}
			return schema + c.proxies.pl[i].p.IP + ":" + c.proxies.pl[i].p.Port, nil
		}
	}

	return "", fmt.Errorf("%s:all(%d),domain(%s)", "No free agent can be used:", length, domain)
}

//IPLimiter rt
func (c *MemCache) IPLimiter(ip string) bool {
	c.ips.l.Lock()
	defer c.ips.l.Unlock()
	now := time.Now().Unix()
	entry, has := c.ips.list[ip]
	if has {
		if now == entry.lastActive && entry.num > proxyinabox.Config.Sys.RequestLimitPerIP {
			return false
		}
	} else {
		c.ips.list[ip] = &ipActivityEntry{num: 1, lastActive: now}
	}

	if entry.lastActive == now {
		entry.num++
	} else {
		entry.num = 1
		entry.lastActive = now
	}
	return true
}

//HostLimiter rt
func (c *MemCache) HostLimiter(ip, host string) bool {

}

//HasProxy rt
func (c *MemCache) HasProxy(p string) bool {

}

//SaveProxy rt
func (c *MemCache) SaveProxy(p proxyinabox.Proxy) error {

}

//DeleteProxy rt
func (c *MemCache) DeleteProxy(p proxyinabox.Proxy) {

}
