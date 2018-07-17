package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/naiba/proxyinabox"
)

//DomainService mysql domain service
type DomainService struct {
	DB *gorm.DB
}

//GetByName get domain by name
func (ds *DomainService) GetByName(name string) (d proxyinabox.Domain, e error) {
	e = ds.DB.Model(&proxyinabox.Domain{}).First(&d, "name = ?", name).Error
	return
}
