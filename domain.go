package proxyinabox

import (
	"github.com/jinzhu/gorm"
)

//Domain proxy request's domain
type Domain struct {
	gorm.Model
	Name string
}
