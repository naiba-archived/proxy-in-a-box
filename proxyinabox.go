package proxyinabox

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/go-redis/redis"
	// mysql driver for GORM
	_ "github.com/go-sql-driver/mysql"
)

//DB instance
var DB *gorm.DB
var Cache *redis.Client

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
		Db   int
	}
	Sys struct {
		Name              string
		ProxyVerifyWorker int   `mapstructure:"proxy_verify_worker"`
		DomainsPerIP      int   `mapstructure:"domains_per_ip"`
		RequestLimitPerIP int32 `mapstructure:"request_limit_per_ip"`
		VerifyDuration    int   `mapstructure:"verify_duration"`
	}
}

//Config system config
var Config Conf

//Init init system
func Init() {
	if Config.Sys.VerifyDuration <= 5 {
		panic("proxy verify duration (must >5 minute)")
	}

	var err error
	DB, err = gorm.Open("mysql", Config.MySQL.User+":"+Config.MySQL.Pass+"@tcp("+Config.MySQL.Host+":"+Config.MySQL.Port+")/"+Config.MySQL.Dbname+"?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("DB!!!", err.Error())
		panic("failed to connect database")
	}

	if Config.Debug {
		DB.LogMode(true)
	}

	DB.AutoMigrate(&Proxy{})

	loadCache()
	clearCacheWorker()
}

func loadCache() {
	var ps []Proxy
	DB.Model(&Proxy{}).Find(&ps)
	proxyQueue.mu.Lock()
	defer proxyQueue.mu.Unlock()
	for _, p := range ps {
		tmp := p
		proxyQueue.list = append(proxyQueue.list, p.ID)
		proxyCache.Store(p.ID, &tmp)
		proxyIndex.Store(p.URI(), p.ID)
	}
}
