package crawler

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/naiba/proxyinabox"
	"github.com/naiba/proxyinabox/service/sqlite3"

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

func init() {
	validateJobs = make(chan proxyinabox.Proxy, 100)
	//start worker
	for i := 1; i <= proxyinabox.ProxyValidatorWorkerNum; i++ {
		go validator(i, validateJobs)
	}
}

func getDocFromURL(url string) (*goquery.Document, error) {

	_, body, errs := gorequest.New().Get(url).
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

//FetchProxies fetch new proxies
func FetchProxies() {
	// in-memory db
	db, err := gorm.Open("sqlite3", "file:box.db?cache=shared&mode=memory&_loc=Asia/Shanghai")
	if err != nil {
		fmt.Println("DB!!!", err.Error())
		panic("failed to connect database")
	}
	db.AutoMigrate(&proxyinabox.Proxy{})
	proxyServiceInstance = &sqlite3.ProxyService{DB: db}

	cs := []proxyinabox.ProxyCrawler{
		newKuai(),
		newXici(),
		new66IP(),
	}

	for _, c := range cs {
		go c.Fetch()
	}
}

func validator(id int, validateJobs chan proxyinabox.Proxy) {
	for p := range validateJobs {
		// format
		p.IP = strings.TrimSpace(p.IP)
		var proxy string
		proxy = "http://" + p.IP + ":" + p.Port
		fmt.Println("worker", id, "process", proxy)
		// 是否正在处理
		_, has := pendingValidate.Load(proxy)
		_, err := proxyServiceInstance.GetByIP(p.IP)
		if !has && err != nil {
			pendingValidate.Store(proxy, nil)
			var resp validateJSON
			start := time.Now().Unix()
			// detect HTTP or HTTPS
			_, _, errs := gorequest.New().Timeout(time.Second*7).Retry(3, time.Second*2, http.StatusInternalServerError).Proxy(proxy).Get("https://api.ip.la/cn?json").EndStruct(&resp)
			if len(errs) != 0 || resp.IP != p.IP {
				start = time.Now().Unix()
				_, _, errs = gorequest.New().Timeout(time.Second*7).Retry(3, time.Second*2, http.StatusInternalServerError).Proxy(proxy).Get("http://api.ip.la/cn?json").EndStruct(&resp)
				p.NotHTTPS = true
			}
			if len(errs) == 0 && resp.IP == p.IP {
				p.Country = resp.Location.CountryName
				p.Provence = resp.Location.Province
				p.Delay = time.Now().Unix() - start

				proxyServiceInstance.Save(&p)
				fmt.Println("worker", id, "find a avaliable proxy", p)
			}
			pendingValidate.Delete(proxy)
		}
	}
}
