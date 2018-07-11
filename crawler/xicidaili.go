package crawler

import (
	"github.com/naiba/proxyinabox"
)

//XiciDaili 西祠代理
type XiciDaili struct {
}

//NewXiciDaili 新建一个西祠代理对象
func NewXiciDaili() *XiciDaili {
	return new(XiciDaili)
}

//GetPage 获取一页中的所有代理
func (cx XiciDaili) GetPage(pageNo int) (list []proxyinabox.Proxy, nextPage int, err error) {
	return nil, 0, nil
}
