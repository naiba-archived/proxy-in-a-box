package sqlite3

import (
	"github.com/jinzhu/gorm"
	"github.com/naiba/proxyinabox"
)

//ProxyService sqlite3 proxy service
type ProxyService struct {
	db *gorm.DB
}

//GetByIP get proxy by ip
func (ps *ProxyService) GetByIP(ip string) (proxyinabox.Proxy, error) {
	var p proxyinabox.Proxy
	return p, ps.db.First(&p, "ip = ?", ip).Error
}

//Save save proxy
func (ps *ProxyService) Save(p proxyinabox.Proxy) error {
	return ps.db.Save(p).Error
}
