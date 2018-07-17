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

//VerifyDuration proxy verify duration (must >5 minute)
const VerifyDuration = 30

//DB instance
var DB *gorm.DB

func init() {

	if VerifyDuration <= 5 {
		panic("proxy verify duration (must >5 minute)")
	}

	// in-memory db "mode=memory"
	var err error
	DB, err = gorm.Open("sqlite3", "file:box.db?cache=shared&mode=rwc&_loc=Asia/Shanghai")
	if err != nil {
		fmt.Println("DB!!!", err.Error())
		panic("failed to connect database")
	}
	//DB = DB.Debug()
	DB.AutoMigrate(&Proxy{}, &Activity{}, &Domain{})
}
