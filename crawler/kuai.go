package crawler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/naiba/proxyinabox"
)

//Kuai 快代理
type Kuai struct {
	urls []string
}

func newKuai() *Kuai {
	this := new(Kuai)
	this.urls = []string{
		"https://www.kuaidaili.com/free/inha/",
		"https://www.kuaidaili.com/free/intr/",
	}
	return this
}

//Fetch fetch all proxies
func (k *Kuai) Fetch() error {

	var currPageNo = 1
	var ended bool

	for _, pageURL := range k.urls {

		for !ended {

			num := 0

			doc, err := getDocFromURL(pageURL + strconv.Itoa(currPageNo))
			if err != nil {
				fmt.Println("[PIAB]", "kuai", "[❎]", "crawler", err)
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
				p.Platform = 2
				//p.HTTPS = strings.Contains(tr.Text(), "HTTPS")
				validateJobs <- p
				num++
			})

			flag := false
			nav := doc.Find("div#listnav").First()
			nav.Find("li").EachWithBreak(func(i int, li *goquery.Selection) bool {
				if strings.TrimSpace(li.Text()) == strconv.Itoa(currPageNo) {
					flag = true
					return true
				}
				if flag {

					currPageNo++

					// 如果当前类型代理遍历完毕
					if strings.TrimSpace(li.Text()) == "页" {
						currPageNo = 1
						ended = true
					}
					return false
				}
				return true
			})

			fmt.Println("[PIAB]", "kuai", "[🍾]", "crawler", num, "proxies.")

			//delay
			time.Sleep(time.Second * 3)
		}
		ended = false
	}

	return nil
}
