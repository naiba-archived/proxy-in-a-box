package proxyinabox

import (
	"github.com/jinzhu/gorm"
)

//Proxy proxy model
type Proxy struct {
	gorm.Model
	IP          string `gorm:"type:varchar(15);unique_index"`
	Port        string `gorm:"type:varchar(5)"`
	Country     string `gorm:"type:varchar(15)"`
	Provence    string `gorm:"type:varchar(15)"`
	IsAnonymous bool
	IsHTTPS     bool
	IsSocks45   bool
	Delay       int64
}

//ProxyCrawler proxy crawler
type ProxyCrawler interface {
	Get() error
}

//ProxyService proxy service
type ProxyService interface {
	GetByIP(ip string) (Proxy, error)
	Save(p Proxy) error
}
