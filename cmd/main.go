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
	c = crawler.NewXici()
	fmt.Println(c.GetPage())
}
