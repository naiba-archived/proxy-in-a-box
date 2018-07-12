package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/naiba/proxyinabox"
	"github.com/naiba/proxyinabox/crawler"
	"github.com/naiba/proxyinabox/service/sqlite3"
)

var ps proxyinabox.ProxyService

func init() {
	db, err := gorm.Open("sqlite3", "box.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.AutoMigrate(&proxyinabox.Proxy{})

	ps = sqlite3.ProxyService{db: db}
	crawler.SetProxyServiceInstance(ps)
}

func main() {
	fmt.Println("AppName:", proxyinabox.AppName)
	fmt.Println("AppVersion:", proxyinabox.AppVersion)

	cs := []proxyinabox.ProxyCrawler{
		crawler.NewKuai(),
		crawler.NewXici(),
	}

	for i := 0; i < 10; i++ {
		for _, c := range cs {
			c.Get()
			time.Sleep(time.Second * 2)
		}
	}

	select {}
}
