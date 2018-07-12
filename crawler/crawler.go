package crawler

import (
	"fmt"
	"strings"

	"github.com/naiba/proxyinabox"

	"github.com/PuerkitoBio/goquery"
	"github.com/naiba/com"
	"github.com/parnurzeal/gorequest"
)

var validateJobs chan proxyinabox.Proxy

func init() {
	if validateJobs == nil {
		validateJobs = make(chan proxyinabox.Proxy, 100)
	}
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
		fmt.Println("worker", id, p)
	}
}
