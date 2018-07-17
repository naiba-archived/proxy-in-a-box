package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/naiba/com"

	"github.com/naiba/proxyinabox"
	"github.com/naiba/proxyinabox/proxy"
)

func main() {
	fmt.Println("AppName:", proxyinabox.AppName)
	fmt.Println("AppVersion:", proxyinabox.AppVersion)

	//crawler.FetchProxies()
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 1)
		proxyinabox.DB.Save(&proxyinabox.Proxy{
			IP:   com.RandomIP(),
			Port: strconv.Itoa(i),
		})
	}

	proxy.Serv("8080", "8081")

	select {}
}
