package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/naiba/com"

	"github.com/robfig/cron"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
		fmt.Println("[Proxy-in-a-Box]", proxyinabox.Config.Sys.Name, "v1.0.0")

		proxy.Serv(httpProxyPort, httpsProxyPort)

		crawler.FetchProxies()
		//crawler.Verify()

		c := cron.New()
		c.AddFunc("@daily", crawler.FetchProxies)
		c.AddFunc("0 "+strconv.Itoa(proxyinabox.Config.Sys.VerifyDuration)+" * * * *", crawler.Verify)
		c.Start()

		select {}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFilePath, "conf", "c", "./pb.yaml", "config file")
	rootCmd.PersistentFlags().StringVarP(&httpProxyPort, "hp", "p", "8080", "http proxy server port")
	rootCmd.PersistentFlags().StringVarP(&httpsProxyPort, "sp", "s", "8081", "https proxy server port")
	//read config
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFilePath)
	com.PanicIfNotNil(viper.ReadInConfig())
	com.PanicIfNotNil(viper.Unmarshal(&proxyinabox.Config))

	proxyinabox.Init()
	crawler.Init()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
