package sqlite3

import (
	"github.com/jinzhu/gorm"
	"github.com/naiba/proxyinabox"
)

//ActivityService activity service
type ActivityService struct {
	DB *gorm.DB
}

//GetByDomainID get activities by domain id
func (as *ActivityService) GetByDomainID(did int64) (a []proxyinabox.Activity, e error) {
	e = as.DB.Model(&proxyinabox.Activity{}).Find(&a, "domain_id = ?", did).Error
	return
}
