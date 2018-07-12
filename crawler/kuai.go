package crawler

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/naiba/proxyinabox"
	"github.com/parnurzeal/gorequest"
)

//Kuai 快代理
type Kuai struct {
	urls       []string
	currURL    int
	currPageNo int
	ended      bool
	req        *gorequest.SuperAgent
}

//NewKuai 新建对象
func NewKuai() *Kuai {
	this := new(Kuai)
	this.urls = []string{
		"https://www.kuaidaili.com/free/inha/",
		"https://www.kuaidaili.com/free/intr/",
	}
	this.currPageNo = 1
	this.req = gorequest.New()
	return this
}

//Get 获取代理
func (k *Kuai) Get() error {
	// 已遍历完毕
	if k.ended {
		return nil
	}

	doc, err := getDocFromURL(k.req, k.urls[k.currURL]+strconv.Itoa(k.currPageNo))
	if err != nil {
		return err
	}

	ipList := doc.Find("div#list table").First()
	ipList.Find("tr").Each(func(i int, tr *goquery.Selection) {

		if i == 0 {
			return
		}

		var p proxyinabox.Proxy
		tr.Children().EachWithBreak(func(j int, td *goquery.Selection) bool {
			if j > 1 {
				return false
			}
			switch j {
			case 0:
				p.IP = td.Text()
			case 1:
				p.Port = td.Text()
			}
			return true
		})
		validateJobs <- p
	})

	flag := false
	nav := doc.Find("div#listnav").First()
	nav.Find("li").EachWithBreak(func(i int, li *goquery.Selection) bool {
		if strings.TrimSpace(li.Text()) == strconv.Itoa(k.currPageNo) {
			flag = true
			return true
		}
		if flag {
			k.currPageNo++

			// 如果当前类型代理遍历完毕
			if strings.TrimSpace(li.Text()) == "页" {
				k.currPageNo = 1
				k.currURL++
				// 如果所有类型代理遍历完毕
				if k.currURL == len(k.urls) {
					k.ended = true
					k.currURL = 0
				}
			}
			return false
		}
		return true
	})
	return nil
}
