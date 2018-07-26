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

func newXici() *Xici {
	this := new(Xici)
	this.urls = []string{
		"http://www.xicidaili.com/nn/",
		"http://www.xicidaili.com/nt/",
	}
	return this
}

//Fetch fetch all proxies
func (xc *Xici) Fetch() error {

	currPageNo := 1
	var ended bool

	for _, pageURL := range xc.urls {
		for !ended {
			num := 0
			doc, err := getDocFromURL(pageURL + strconv.Itoa(currPageNo))
			if err != nil {
				fmt.Println("[PIAB]", "xici", "[❎]", "crawler", err)
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
				p.Platform = 3
				//p.HTTPS = strings.Contains(tr.Text(), "HTTPS")

				validateJobs <- p
				num++
			})

			currPageNo++

			nextPage := doc.Find("span.next_page").First()
			// 如果当前类型代理遍历完毕
			if nextPage.HasClass("disabled") {
				currPageNo = 1
				ended = true
			}

			fmt.Println("[PIAB]", "xici", "[🍾]", "crawler", num, "proxies.")

			//delay
			time.Sleep(time.Second * 3)
		}
		ended = false
	}
	return nil
}
