package sqlite3

import (
	"github.com/jinzhu/gorm"
	"github.com/naiba/proxyinabox"
)

//DomainService sqlite3 domain service
type DomainService struct {
	DB *gorm.DB
}

//GetByName get domain by name
func (ds *DomainService) GetByName(name string) (d proxyinabox.Domain, e error) {
	e = ds.DB.Model(&proxyinabox.Domain{}).First(&d, "name = ?", name).Error
	return
}
