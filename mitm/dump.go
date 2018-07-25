package mitm

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

//Dump rt
func (m *MITM) Dump(resp http.ResponseWriter, req *http.Request, https bool) {
	if m.IsDirect {
		req.Header.Del("Proxy-Connection")
		req.Header.Set("Connection", "Keep-Alive")
	}
	var reqDump []byte
	var err error
	ch := make(chan bool)
	go func() {
		reqDump, err = httputil.DumpRequestOut(req, true)
		if err != nil {
			fmt.Println("DumpRequest error ", err)
		}
		ch <- true
	}()
	connIn, _, err := resp.(http.Hijacker).Hijack()
	if err != nil {
		fmt.Println("hijack error:", err)
	}
	defer connIn.Close()

	var respOut *http.Response
	host := req.Host

	var connOut net.Conn
	if !https {
		connOut, err = net.DialTimeout("tcp", host, time.Second*30)
	} else {
		connOut, err = tls.Dial("tcp", host, m.TLSConf.ServerTLSConfig)
	}
	if err != nil {
		fmt.Println("tls dial to", host, "error:", err)
		return
	}
	defer connOut.Close()

	if err = req.Write(connOut); err != nil {
		fmt.Println("send to server error", err)
		return
	}

	respOut, err = http.ReadResponse(bufio.NewReader(connOut), req)
	if err != nil && err != io.EOF {
		fmt.Println("read response error:", err)
		return
	}
	if respOut == nil {
		fmt.Println("respOut is nil")
		return
	}

	respDump, err := httputil.DumpResponse(respOut, true)
	if err != nil {
		fmt.Println("respDump error:", err)
		return
	}

	_, err = connIn.Write(respDump)
	if err != nil {
		fmt.Println("connIn write error:", err)
		return
	}

	fmt.Println("REQUEST:", string(reqDump))
	fmt.Println("RESPONSE:", string(respDump))
	<-ch
}
