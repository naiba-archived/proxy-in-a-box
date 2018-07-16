package main

import (
	"fmt"

	"github.com/naiba/proxyinabox"
	"github.com/naiba/proxyinabox/proxy"
)

func main() {
	fmt.Println("AppName:", proxyinabox.AppName)
	fmt.Println("AppVersion:", proxyinabox.AppVersion)

	//crawler.FetchProxies()

	proxy.Serv("8080", "8081")

	select {}
}
