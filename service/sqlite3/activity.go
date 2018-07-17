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
func (as *ActivityService) GetByDomainID(did uint) (a []proxyinabox.Activity, e error) {
	e = as.DB.Find(&a, "domain_id = ?", did).Error
	return
}

//Save save activity
func (as *ActivityService) Save(d, p uint) {
	var a proxyinabox.Activity
	if e := as.DB.First(&a, "domain_id = ? AND proxy_id = ?", d, p).Error; e == nil {
		a.Usenum++
		as.DB.Save(&a)
	} else if e == gorm.ErrRecordNotFound {
		as.DB.Save(&proxyinabox.Activity{
			DomainID: d,
			ProxyID:  p,
		})
	}
}
