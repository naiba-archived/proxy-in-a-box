package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/robfig/cron"
	"github.com/spf13/cobra"

	"github.com/naiba/proxyinabox"
	"github.com/naiba/proxyinabox/crawler"
	"github.com/naiba/proxyinabox/proxy"
)

var configFilePath, httpProxyPort, httpsProxyPort string
var rootCmd = &cobra.Command{
	Use:   "proxy-in-a-box",
	Short: "Proxy-in-a-Box provide many proxies.",
	Long:  `Proxy-in-a-Box helps programmers quickly and easily develop powerful crawler services. one-script, easy-to-use: proxies in a box.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("AppName:", proxyinabox.AppName)
		fmt.Println("AppVersion:", proxyinabox.AppVersion)

		crawler.FetchProxies()
		//crawler.Verify()

		c := cron.New()
		c.AddFunc("@daily", crawler.FetchProxies)
		c.AddFunc("0 "+strconv.Itoa(proxyinabox.VerifyDuration)+" * * * *", crawler.Verify)
		c.Start()

		proxy.Serv(httpProxyPort, httpsProxyPort)

		fmt.Println("HTTP proxy: `http://localhost:" + httpProxyPort + "`")
		fmt.Println("HTTPS proxy: `https://localhost:" + httpsProxyPort + "`")

		select {}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFilePath, "conf", "c", "./pb.yaml", "config file")
	rootCmd.PersistentFlags().StringVarP(&httpProxyPort, "hp", "p", "8080", "http proxy server port")
	rootCmd.PersistentFlags().StringVarP(&httpsProxyPort, "sp", "s", "8081", "https proxy server port")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
