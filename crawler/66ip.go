package crawler

import (
	"fmt"
	"strings"
	"time"

	"github.com/naiba/com"
	"github.com/naiba/proxyinabox"
	"github.com/parnurzeal/gorequest"
)

//P66IP 66ip site
type P66IP struct {
	urls []string
}

func new66IP() *P66IP {
	this := new(P66IP)
	this.urls = []string{
		"http://www.66ip.cn/mo.php?tqsl=200",
		"http://www.66ip.cn/nmtq.php?getnum=200",
	}
	return this
}

//Fetch fetch all proxies
func (p *P66IP) Fetch() error {
	for _, pageURL := range p.urls {
		for i := 0; i < 50; i++ {
			num := 0
			_, body, errs := gorequest.New().Get(pageURL).
				Set("User-Agent", com.RandomUserAgent()).
				Retry(3, time.Second*3).
				End()
			if len(errs) != 0 {
				fmt.Println("[PIAB]", "66ip", "[❎]", "crawler", errs)
				continue
			}
			lines := strings.Split(string(body), "<br />")
			for _, line := range lines {
				ipinfo := strings.Split(strings.TrimSpace(line), ":")
				if len(ipinfo) == 2 && com.IsIPv4(ipinfo[0]) {
					var p proxyinabox.Proxy
					p.IP = ipinfo[0]
					p.Port = ipinfo[1]
					p.Platform = 1
					p.HTTPS = false
					num++
					validateJobs <- p
				}
			}
			fmt.Println("[PIAB]", "66ip", "[🍾]", "crawler", num, "proxies.")
			//delay
			time.Sleep(time.Second * 3)
		}
	}
	return nil
}
