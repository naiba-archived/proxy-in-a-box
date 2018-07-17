package main

import (
	"fmt"
	"strconv"

	"github.com/robfig/cron"

	"github.com/naiba/proxyinabox"
	"github.com/naiba/proxyinabox/crawler"
	"github.com/naiba/proxyinabox/proxy"
)

func main() {
	fmt.Println("AppName:", proxyinabox.AppName)
	fmt.Println("AppVersion:", proxyinabox.AppVersion)

	c := cron.New()
	c.AddFunc("@daily", crawler.FetchProxies)
	c.AddFunc("0 "+strconv.Itoa(proxyinabox.VerifyDuration)+" * * * *", crawler.Verify)
	c.Start()

	proxy.Serv("8080", "8081")

	select {}
}
