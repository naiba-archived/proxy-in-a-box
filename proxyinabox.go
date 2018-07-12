package proxyinabox

import (
	"strings"

	"github.com/parnurzeal/gorequest"
)

//AppName app's name
const AppName = "Proxy-in-a-Box"

//AppVersion app's version
const AppVersion = "1.0"

//ProxyValidatorWorkerNum 代理验证worker数
const ProxyValidatorWorkerNum = 20

//ServerIP app's ip address
var ServerIP string

func init() {
	_, body, err := gorequest.New().Get("https://api.ip.la").End()
	if err != nil {
		panic(err)
	}
	ServerIP = strings.TrimSpace(body)
}
