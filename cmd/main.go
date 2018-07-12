package main

import (
	"fmt"
	"time"

	"github.com/naiba/proxyinabox"
	"github.com/naiba/proxyinabox/crawler"
)

func main() {
	fmt.Println("AppName:", proxyinabox.AppName)
	fmt.Println("AppVersion:", proxyinabox.AppVersion)

	cs := []proxyinabox.ProxyCrawler{
		crawler.NewKuai(),
		crawler.NewXici(),
	}

	for i := 0; i < 2; i++ {
		for _, c := range cs {
			c.Get()
			time.Sleep(time.Second * 2)
		}
	}

	select {}
}
