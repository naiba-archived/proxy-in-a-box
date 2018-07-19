package proxyinabox

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

//Proxy proxy model
type Proxy struct {
	gorm.Model
	IP         string `gorm:"type:varchar(15);unique_index"`
	Port       string `gorm:"type:varchar(5)"`
	Country    string `gorm:"type:varchar(15)"`
	Provence   string `gorm:"type:varchar(15)"`
	Platform   int
	NotHTTPS   bool
	Delay      int64
	LastVerify time.Time
}

//ProxyService proxy service
type ProxyService interface {
	GetUnVerified() ([]Proxy, error)
}

func (p Proxy) String() string {
	return fmt.Sprintf("[√]proxy#[id:%d %s:%s country:%s provence:%s HTTPS:%t delay:%d platform:%d]",
		p.ID, p.IP, p.Port, p.Country, p.Provence, !p.NotHTTPS, p.Delay, p.Platform)
}

//URI get uri
func (p Proxy) URI() string {
	return "http://" + p.IP + ":" + p.Port
}

//ProxyCrawler proxy crawler
type ProxyCrawler interface {
	Fetch() error
}
