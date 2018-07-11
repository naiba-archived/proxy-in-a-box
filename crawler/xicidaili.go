package crawler

import (
	"log"
	"strconv"
	"strings"

	"github.com/naiba/com"

	"github.com/PuerkitoBio/goquery"
	"github.com/naiba/proxyinabox"
	"github.com/parnurzeal/gorequest"
)

//XiciDaili 西祠代理
type XiciDaili struct {
	urls    []string
	currURL int
}

//NewXiciDaili 新建一个西祠代理对象
func NewXiciDaili() *XiciDaili {
	this := new(XiciDaili)
	this.urls = []string{
		"http://www.xicidaili.com/nn/",
		"http://www.xicidaili.com/nt/",
		"http://www.xicidaili.com/qq/",
	}
	return this
}

//GetPage 获取一页中的所有代理
func (xc XiciDaili) GetPage(pageNo int) (list []proxyinabox.Proxy, nextPageNo int, err error) {
	if pageNo == 0 {
		pageNo = 1
	}

	request := gorequest.New()
	_, body, _ := request.Get(xc.urls[xc.currURL]+strconv.Itoa(pageNo)).
		Set("User-Agent", com.RandomUserAgent()).
		End()

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(body))

	ipList := doc.Find("table#ip_list").First()
	ipList.Find("tr").Each(func(i int, tr *goquery.Selection) {

		if i == 0 {
			return
		}

		var ip string
		var port string
		tr.Children().EachWithBreak(func(j int, td *goquery.Selection) bool {
			if j > 2 {
				return false
			}
			switch j {
			case 1:
				ip = td.Text()
			case 2:
				port = td.Text()
			}
			return true
		})

		log.Println("IP", ip, "Port", port)
	})
	return
}
