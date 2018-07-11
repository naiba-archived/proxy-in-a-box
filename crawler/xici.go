package crawler

import (
	"strconv"
	"strings"

	"github.com/naiba/com"

	"github.com/PuerkitoBio/goquery"
	"github.com/naiba/proxyinabox"
	"github.com/parnurzeal/gorequest"
)

//Xici 西祠代理
type Xici struct {
	urls       []string
	currURL    int
	currPageNo int
	ended      bool
	req        *gorequest.SuperAgent
}

//NewXici 新建一个西祠代理对象
func NewXici() *Xici {
	this := new(Xici)
	this.urls = []string{
		"http://www.xicidaili.com/nn/",
		"http://www.xicidaili.com/nt/",
	}
	this.currPageNo = 1
	this.req = gorequest.New()
	return this
}

//Get 获取一页代理，会自动翻页、换类型
func (xc *Xici) Get() (list []proxyinabox.Proxy, err error) {

	// 已遍历完毕
	if xc.ended {
		return
	}

	_, body, errs := xc.req.Get(xc.urls[xc.currURL]+strconv.Itoa(xc.currPageNo)).
		Set("User-Agent", com.RandomUserAgent()).
		End()
	if len(errs) > 0 {
		err = errs[0]
		return
	}

	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return
	}

	list = make([]proxyinabox.Proxy, 0)

	ipList := doc.Find("table#ip_list").First()
	ipList.Find("tr").Each(func(i int, tr *goquery.Selection) {

		if i == 0 {
			return
		}

		var p proxyinabox.Proxy
		tr.Children().EachWithBreak(func(j int, td *goquery.Selection) bool {
			if j > 2 {
				return false
			}
			switch j {
			case 1:
				p.IP = td.Text()
			case 2:
				p.Port = td.Text()
			}
			return true
		})

		content := tr.Text()
		p.IsAnonymous = strings.Contains(content, "高匿")
		p.IsHTTPS = strings.Contains(content, "HTTPS")
		p.IsSocks45 = strings.Contains(content, "socks4/5")

		list = append(list, p)
	})

	xc.currPageNo++

	nextPage := doc.Find("span.next_page").First()
	// 如果当前类型代理遍历完毕
	if nextPage.HasClass("disabled") {
		xc.currPageNo = 1
		xc.currURL++

		// 如果所有类型代理遍历完毕
		if xc.currURL == len(xc.urls) {
			xc.ended = true
			xc.currURL = 0
		}
	}
	return
}
