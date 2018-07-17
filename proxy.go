package proxyinabox

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

//Proxy proxy model
type Proxy struct {
	gorm.Model
	IP       string `gorm:"type:varchar(15);unique_index"`
	Port     string `gorm:"type:varchar(5)"`
	Country  string `gorm:"type:varchar(15)"`
	Provence string `gorm:"type:varchar(15)"`
	Platform int
	NotHTTPS bool
	Delay    int64
	Usenum   int64

	Activitys []Activity
}

func (p Proxy) String() string {
	return fmt.Sprintf("[√]proxy#[id:%d %s:%s country:%s provence:%s HTTPS:%t delay:%d useNum:%d platform:%d]",
		p.ID, p.IP, p.Port, p.Country, p.Provence, !p.NotHTTPS, p.Delay, p.Usenum, p.Platform)
}

//ProxyCrawler proxy crawler
type ProxyCrawler interface {
	Fetch() error
}

//ProxyService proxy service
type ProxyService interface {
	GetByIP(ip string) (Proxy, error)
	GetFree(notIn []uint) (Proxy, error)
}
