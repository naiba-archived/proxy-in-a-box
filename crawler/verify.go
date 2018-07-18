package crawler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/naiba/proxyinabox"
	"github.com/parnurzeal/gorequest"
)

var verifyJob chan proxyinabox.Proxy

func init() {
	verifyJob = make(chan proxyinabox.Proxy, proxyinabox.Config.Sys.ProxyVerifyWorker)
	for i := 0; i < proxyinabox.Config.Sys.ProxyVerifyWorker; i++ {
		go getDelay(verifyJob)
	}
}

//Verify verify proxies in database
func Verify() {
	list, err := proxyServiceInstance.GetUnVerified()
	if err != nil {
		fmt.Println(err)
		Verify()
		return
	}
	for _, p := range list {
		verifyJob <- p
	}
}

func getDelay(pc chan proxyinabox.Proxy) {
	for p := range pc {
		proxy := "http://" + p.IP + ":" + p.Port
		start := time.Now().Unix()
		var resp validateJSON
		_, _, errs := gorequest.New().Timeout(time.Second*5).Retry(3, time.Second*2, http.StatusInternalServerError).Proxy(proxy).Get("http://api.ip.la/cn?json").EndStruct(&resp)
		delay := time.Now().Unix() - start
		if len(errs) != 0 || resp.IP != p.IP {
			proxyinabox.DB.Delete(&p)
			continue
		}
		p.Delay = delay
		proxyinabox.DB.Update(&p)
	}
}
