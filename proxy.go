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
}

//ProxyCrawler 代理抓取工具
type ProxyCrawler interface {
	GetPage(pageNo int) (list []Proxy, nextPageNo int, err error)
}
