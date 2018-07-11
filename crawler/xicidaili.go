package crawler

import (
	"github.com/naiba/proxyinabox"
)

//XiciDaili 西祠代理
type XiciDaili struct {
	urls    []string
	currURL int
}

//NewXiciDaili 新建一个西祠代理对象
func NewXiciDaili() *XiciDaili {
	this := new(XiciDaili)
	this.urls = []string{
		"http://www.xicidaili.com/nn/",
		"http://www.xicidaili.com/nt/",
		"http://www.xicidaili.com/qq",
	}
	return this
}

//GetPage 获取一页中的所有代理
func (cx XiciDaili) GetPage(pageNo int) (list []proxyinabox.Proxy, nextPageNo int, err error) {
	if pageNo == 0 {
		pageNo = 1
	}

	return nil, 0, nil
}
