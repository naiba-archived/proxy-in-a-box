package mitm

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//Dump rt
func (m *MITM) Dump(clientResponse http.ResponseWriter, clientRequest *http.Request) {
	var clientRequestDump []byte
	var remoteResponseDump []byte
	var remoteResponse *http.Response
	var err error
	defer func() {
		if err != nil {
			clientResponse.WriteHeader(http.StatusBadGateway)
			clientResponse.Write([]byte(err.Error()))
		}
	}()
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
		fmt.Println("[MITM]", "remoteResponse", "[❎]", err)
		return
	}

	remoteResponseDump, err = httputil.DumpResponse(remoteResponse, true)
	if err != nil {
		fmt.Println("[MITM]", "respDump", "[❎]", err)
		return
	}
	// set response header
	for k, v := range remoteResponse.Header {
		if k == "Content-Encoding" {
			continue
		}
		var vb []byte
		for i := 0; i < len(v); i++ {
			if i == len(v)-1 {
				vb = append(vb, []byte(v[i])...)
			} else {
				vb = append(vb, []byte(v[i]+"; ")...)
			}
		}
		clientResponse.Header().Set(k, string(vb))
	}
	// write response code
	clientResponse.WriteHeader(remoteResponse.StatusCode)
	// decompress gzip page
	body := make([]byte, 0)
	switch remoteResponse.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ := gzip.NewReader(remoteResponse.Body)
		var n int
		for {
			buf := make([]byte, 1024)
			n, err = reader.Read(buf)

			if err != nil && err != io.EOF {
				fmt.Println("[MITM]", "decompress gzip", "[❎]", err)
				break
			}

			if n == 0 {
				break
			}
			body = append(body, buf...)
		}
	default:
		body, err = ioutil.ReadAll(remoteResponse.Body)
	}
	// write response body
	if err != nil {
		fmt.Println("[MITM]", "read body", "[❎]", err)
		return
	}
	_, err = clientResponse.Write(body)
	if err != nil {
		fmt.Println("[MITM]", "connIn write", "[❎]", err)
		return
	}
	// show http dump
	if m.Print {
		fmt.Println("[MITM]", "REQUEST-DUMP", "[📮]", string(clientRequestDump))
		fmt.Println("[MITM]", "RESPONSE-DUMP", "[📮]", string(remoteResponseDump))

	}
	<-ch
}
