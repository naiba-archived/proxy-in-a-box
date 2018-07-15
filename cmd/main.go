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
	// in-memory db
	db, err := gorm.Open("sqlite3", "file:box.db?cache=shared&mode=memory&_loc=Asia/Shanghai")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&proxyinabox.Proxy{})

	ps = &sqlite3.ProxyService{DB: db}
	crawler.SetProxyServiceInstance(ps)
}

func main() {
	fmt.Println("AppName:", proxyinabox.AppName)
	fmt.Println("AppVersion:", proxyinabox.AppVersion)

	// TODO: one trigger crawl all pages
	cs := []proxyinabox.ProxyCrawler{
		crawler.NewKuai(),
		crawler.NewXici(),
		crawler.New66IP(),
	}

	for i := 0; i < 100; i++ {
		for _, c := range cs {
			c.Get()
			time.Sleep(time.Second * 2)
		}
	}

	select {}
}
