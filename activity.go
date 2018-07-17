package proxyinabox

import (
	"github.com/jinzhu/gorm"
)

//Activity proxy use activity
type Activity struct {
	gorm.Model
	Domain   Domain
	DomainID int64
	Proxy    Proxy
	ProxyID  int64
}

//ActivityService activity service
type ActivityService interface {
	GetByDomainID(did int64) ([]Activity, error)
}
