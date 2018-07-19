package proxy

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/naiba/proxyinabox"
)

func handleTunneling(proxy *proxyinabox.Proxy, w http.ResponseWriter, r *http.Request) {
	//set proxy
	destConn, err := net.DialTimeout("tcp", proxy.IP+":"+proxy.Port, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	//send connect message
	destConn.Write([]byte("CONNECT " + r.Host + " HTTP/1.1\r\n\r\n"))
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}
func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
func handleHTTP(proxy *proxyinabox.Proxy, w http.ResponseWriter, req *http.Request) {
	//set proxy
	p, _ := url.Parse("http://" + proxy.IP + ":" + proxy.Port)
	tp := &http.Transport{
		Proxy: http.ProxyURL(p),
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	resp, err := tp.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
