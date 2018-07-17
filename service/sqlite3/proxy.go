package sqlite3

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/naiba/proxyinabox"
)

//ProxyService sqlite3 proxy service
type ProxyService struct {
	DB *gorm.DB
}

//GetByIP get proxy by ip
func (ps *ProxyService) GetByIP(ip string) (proxyinabox.Proxy, error) {
	var p proxyinabox.Proxy
	return p, ps.DB.First(&p, "ip = ?", ip).Error
}

//GetFree get a free proxy
func (ps *ProxyService) GetFree(notIn []uint) (p proxyinabox.Proxy, e error) {
	e = ps.DB.Not(notIn).Order("usenum ASC,delay ASC").First(&p).Error
	return
}

//GetUnVerified get un verified proxies
func (ps *ProxyService) GetUnVerified() (p []proxyinabox.Proxy, e error) {
	e = ps.DB.Where("last_verify < ?", time.Now().Add(time.Minute*(proxyinabox.VerifyDuration-5)*-1)).Find(&p).Error
	return
}
