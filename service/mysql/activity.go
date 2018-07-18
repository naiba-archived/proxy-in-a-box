package mysql

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
	// save avtivity
	var a proxyinabox.Activity
	if e := as.DB.Select("id,usenum").First(&a, "domain_id = ? AND proxy_id = ?", d, p).Error; e == nil {
		a.Usenum++
		as.DB.Model(&a).Update("usenum", a.Usenum)
	} else if e == gorm.ErrRecordNotFound {
		as.DB.Save(&proxyinabox.Activity{
			DomainID: d,
			ProxyID:  p,
			Usenum:   1,
		})
	}
	// update proxy use
	var proxy proxyinabox.Proxy
	if as.DB.Select("id,usenum").First(&proxy, p).Error == nil {
		proxy.Usenum++
		as.DB.Model(&proxy).Update("usenum", proxy.Usenum)
	}
}
