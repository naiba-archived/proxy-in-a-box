package crawler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/naiba/proxyinabox"
)

//Xici 西祠代理
type Xici struct {
	urls []string
}

//NewXici 新建一个西祠代理对象
func NewXici() *Xici {
	this := new(Xici)
	this.urls = []string{
		"http://www.xicidaili.com/nn/",
		"http://www.xicidaili.com/nt/",
	}
	return this
}

//Get 获取一页代理，会自动翻页、换类型
func (xc *Xici) Get() error {

	currPageNo := 1
	var ended bool

	for _, pageURL := range xc.urls {
		for !ended {
			doc, err := getDocFromURL(pageURL + strconv.Itoa(currPageNo))
			if err != nil {
				fmt.Println("XICI ERROR!!", err.Error())
				return err
			}

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
				validateJobs <- p
			})

			currPageNo++

			nextPage := doc.Find("span.next_page").First()
			// 如果当前类型代理遍历完毕
			if nextPage.HasClass("disabled") {
				currPageNo = 1
				ended = true
			}

			//delay
			time.Sleep(time.Second * 3)
		}
		ended = false
	}
	return nil
}
