package proxyinabox

import (
	"github.com/jinzhu/gorm"
)

//Activity 活动列表
type Activity struct {
	gorm.Model
	DomainID int64
}
