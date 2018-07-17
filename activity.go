package proxyinabox

import (
	"github.com/jinzhu/gorm"
)

//Activity proxy use activity
type Activity struct {
	gorm.Model
	Domain   Domain
	DomainID uint `gorm:"index"`
	Proxy    Proxy
	ProxyID  uint `gorm:"index"`
	Usenum   int64
}

//ActivityService activity service
type ActivityService interface {
	GetByDomainID(did uint) ([]Activity, error)
	Save(d, p uint)
}