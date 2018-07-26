package proxyinabox

import (
	"net/http"
)

//Cache rt
type Cache interface {
	PickProxy(req *http.Request) (string, error)
	IPLimiter(ip string) bool
	HostLimiter(ip, host string) bool
	HasProxy(p string) bool
	SaveProxy(p Proxy) error
	DeleteProxy(p Proxy)
}
