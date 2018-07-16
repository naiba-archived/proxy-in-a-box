package proxyinabox

import (
	"github.com/jinzhu/gorm"
)

//Activity 活动列表
type Activity struct {
	gorm.Model
	Domain   Domain
	DomainID int64
	Proxy    Proxy
	ProxyID  int64
}
