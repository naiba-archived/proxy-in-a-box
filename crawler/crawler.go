package crawler

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/naiba/proxyinabox"

	"github.com/PuerkitoBio/goquery"
	"github.com/naiba/com"
	"github.com/parnurzeal/gorequest"
)

var validateJobs chan proxyinabox.Proxy
var pendingValidate sync.Map
var proxyServiceInstance proxyinabox.ProxyService

type validateJSON struct {
	IP       string
	Location struct {
		City        string
		CountryCode string `json:"country_code"`
		CountryName string `json:"country_name"`
		Latitude    string
		Longitude   string
		Province    string
	}
}

//SetProxyServiceInstance set proxy service instance
func SetProxyServiceInstance(ps proxyinabox.ProxyService) {
	proxyServiceInstance = ps
}

func init() {
	validateJobs = make(chan proxyinabox.Proxy, 100)
	//start worker
	for i := 1; i <= proxyinabox.ProxyValidatorWorkerNum; i++ {
		go validator(i, validateJobs)
	}
}

func getDocFromURL(req *gorequest.SuperAgent, url string) (*goquery.Document, error) {

	_, body, errs := req.Get(url).
		Set("User-Agent", com.RandomUserAgent()).
		End()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	return doc, nil

}

func validator(id int, validateJobs chan proxyinabox.Proxy) {
	for p := range validateJobs {
		var proxy string
		if p.IsSocks45 {
			proxy = "socks5://" + p.IP + ":" + p.Port
		} else {
			proxy = "http://" + p.IP + ":" + p.Port
		}
		// 是否正在处理
		_, has := pendingValidate.Load(proxy)
		_, err := proxyServiceInstance.GetByIP(p.IP)
		if !has && err != nil {
			pendingValidate.Store(proxy, nil)
			var resp validateJSON
			var ipip string

			if p.IsSocks45 || p.IsHTTPS {
				ipip = "https://api.ip.la/cn?json"
			} else {
				ipip = "http://api.ip.la/cn?json"
			}
			start := time.Now().Unix()
			_, _, errs := gorequest.New().Timeout(time.Second*7).Retry(3, time.Second*2, http.StatusInternalServerError).Proxy(proxy).Get(ipip).EndStruct(&resp)
			if len(errs) == 0 && resp.IP == p.IP {
				p.Country = resp.Location.CountryName
				p.Provence = resp.Location.Province
				p.Delay = time.Now().Unix() - start

				fmt.Println("worker", id, "find a avaliable proxy", proxy)

				proxyServiceInstance.Save(&p)
			}
			pendingValidate.Delete(proxy)
		}
	}
}
