package mitm

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//Dump rt
func (m *MITM) Dump(resp http.ResponseWriter, req *http.Request) {
	var reqDump []byte
	var respDump []byte
	var err error
	var respOut *http.Response
	ch := make(chan bool)
	go func() {
		reqDump, err = httputil.DumpRequestOut(req, true)
		if err != nil {
			fmt.Println("DumpRequest error ", err)
		}
		ch <- true
	}()

	tp := http.Transport{}

	if !m.IsDirect {
		proxy, err := m.Scheduler(req)
		if err != nil {
			fmt.Println("prxy scheduler error", err)
			return
		}
		p, err := url.Parse(proxy)
		if err != nil {
			fmt.Println("prxy parse error", err)
			return
		}
		tp.Proxy = http.ProxyURL(p)
	} else {
		req.Header.Del("Proxy-Connection")
		req.Header.Set("Connection", "Keep-Alive")
	}

	req.RequestURI = ""
	cli := http.Client{Transport: &tp}
	respOut, err = cli.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	respDump, err = httputil.DumpResponse(respOut, true)
	if err != nil {
		fmt.Println("respDump error:", err)
		return
	}

	resp.WriteHeader(respOut.StatusCode)
	_, err = resp.Write(respDump)
	if err != nil {
		fmt.Println("connIn write error:", err)
		return
	}

	fmt.Println("REQUEST:", string(reqDump))
	fmt.Println("RESPONSE:", string(respDump))
	<-ch
}
