package mitm

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//Dump rt
func (m *MITM) Dump(clientResponse http.ResponseWriter, clientRequest *http.Request) {
	var clientRequestDump []byte
	var remoteResponseDump []byte
	var err error
	var remoteResponse *http.Response
	ch := make(chan bool)
	go func() {
		clientRequestDump, err = httputil.DumpRequestOut(clientRequest, true)
		if err != nil {
			fmt.Println("[MITM]", "DumpRequest", "[❎]", err)
		}
		ch <- true
	}()

	transport := http.Transport{}

	if !m.IsDirect {
		proxy, err := m.Scheduler(clientRequest)
		if err != nil {
			fmt.Println("[MITM]", "proxy scheduler", "[❎]", err)
			return
		}
		p, err := url.Parse(proxy)
		if err != nil {
			fmt.Println("[MITM]", "proxy parse", "[❎]", err)
			return
		}
		transport.Proxy = http.ProxyURL(p)
	} else {
		clientRequest.Header.Del("Proxy-Connection")
		clientRequest.Header.Set("Connection", "Keep-Alive")
	}

	clientRequest.RequestURI = ""
	cli := http.Client{
		Transport: &transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return fmt.Errorf("")
		},
	}
	remoteResponse, err = cli.Do(clientRequest)
	if err != nil {
		fmt.Println("[MITM]", "proxy parse", "[❎]", err)
		return
	}

	remoteResponseDump, err = httputil.DumpResponse(remoteResponse, true)
	if err != nil {
		fmt.Println("[MITM]", "respDump", "[❎]", err)
		return
	}

	clientResponse.WriteHeader(remoteResponse.StatusCode)
	_, err = clientResponse.Write(remoteResponseDump)
	if err != nil {
		fmt.Println("[MITM]", "connIn write", "[❎]", err)
		return
	}

	fmt.Println("[MITM]", "REQUEST", "[📮]", string(clientRequestDump))
	fmt.Println("[MITM]", "RESPONSE", "[📮]", string(remoteResponseDump))
	<-ch
}
