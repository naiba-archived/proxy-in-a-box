package proxyinabox

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// mysql driver for GORM
	_ "github.com/go-sql-driver/mysql"
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
	DB, err = gorm.Open("mysql", "root:123456@tcp(localhost:3306)/proxy?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("DB!!!", err.Error())
		panic("failed to connect database")
	}
	// DB = DB.Debug()
	DB.AutoMigrate(&Proxy{}, &Activity{}, &Domain{})
}
