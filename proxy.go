package proxyinabox

//Proxy 代理
type Proxy struct {
	IP          string
	Port        string
	Country     string
	Provence    string
	IsAnonymous bool
	IsHTTPS     bool
	IsSocks45   bool
	Delay       int
}

//ProxyCrawler 代理抓取工具
type ProxyCrawler interface {
	Get() (list []Proxy, err error)
}
