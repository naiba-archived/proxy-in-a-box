package main

import (
	"fmt"

	"github.com/naiba/proxyinabox"
	"github.com/naiba/proxyinabox/crawler"
)

func main() {
	fmt.Println("AppName:", proxyinabox.AppName)
	fmt.Println("AppVersion:", proxyinabox.AppVersion)

	var c proxyinabox.ProxyCrawler
	c = crawler.NewXiciDaili()
	fmt.Println(c.GetPage(0))
}
