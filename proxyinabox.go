package proxyinabox

import (
	"fmt"

	"github.com/jinzhu/gorm"

	// mysql driver for GORM
	_ "github.com/go-sql-driver/mysql"
)

//DB instance
var DB *gorm.DB

//CI cache instance
var CI Cache

//Conf config struct
type Conf struct {
	Debug bool
	MySQL struct {
		Host   string
		Port   string
		User   string
		Pass   string
		Dbname string
	} `mapstructure:"mysql"`
	Redis struct {
		Host string
		Port string
		Pass string
		Db   int
	}
	Sys struct {
		Name              string
		ProxyVerifyWorker int   `mapstructure:"proxy_verify_worker"`
		DomainsPerIP      int   `mapstructure:"domains_per_ip"`
		RequestLimitPerIP int64 `mapstructure:"request_limit_per_ip"`
		VerifyDuration    int   `mapstructure:"verify_duration"`
	}
}

//Config system config
var Config Conf

//Init init system
func Init() {
	validateConf()
	initDB()
	loadCache()
}

func loadCache() {
	var ps []Proxy
	DB.Model(&Proxy{}).Find(&ps)
	for _, p := range ps {
		//TODO:loadCache
		fmt.Println(p)
	}
}

func initDB() {
	var err error
	DB, err = gorm.Open("mysql", Config.MySQL.User+":"+Config.MySQL.Pass+"@tcp("+Config.MySQL.Host+":"+Config.MySQL.Port+")/"+Config.MySQL.Dbname+"?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	if Config.Debug {
		DB.LogMode(true)
	}
	DB.AutoMigrate(&Proxy{})
}

func validateConf() {
	if Config.Sys.VerifyDuration <= 5 {
		panic("proxy verify duration (must >5 minute)")
	}
}
