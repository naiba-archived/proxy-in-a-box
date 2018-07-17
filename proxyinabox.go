package proxyinabox

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// sqlite driver for GORM
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//AppName app's name
const AppName = "Proxy-in-a-Box"

//AppVersion app's version
const AppVersion = "1.0"

//ProxyValidatorWorkerNum verify proxy's worker num
const ProxyValidatorWorkerNum = 20

//DomainsPerIPHalfAnHour domains num per ip on half hour
const DomainsPerIPHalfAnHour = 10

//DB instance
var DB *gorm.DB

func init() {
	// in-memory db
	var err error
	DB, err = gorm.Open("sqlite3", "file:box.db?cache=shared&mode=memory&_loc=Asia/Shanghai")
	if err != nil {
		fmt.Println("DB!!!", err.Error())
		panic("failed to connect database")
	}
	//DB = DB.Debug()
	DB.AutoMigrate(&Proxy{}, &Activity{}, &Domain{})
}
