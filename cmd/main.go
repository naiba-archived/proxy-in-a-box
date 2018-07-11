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

	var c proxyinabox.ProxyCrawler
	c = crawler.NewKuai()
	for i := 0; i < 4; i++ {
		fmt.Println(c.Get())
		time.Sleep(time.Second)
	}
}
