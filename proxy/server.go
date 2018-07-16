package proxy

import (
	"crypto/tls"
	"net/http"
	"strconv"
	"strings"

	"github.com/naiba/proxyinabox"
)

//Serv serv the http proxy
func Serv(httpPort, httpsPort string) {
	httpServer := newServer(httpPort)
	go httpServer.ListenAndServe()

	var pemPath = "../tls_keys/server.pem"
	var keyPath = "../tls_keys/server.key"
	httpsServer := newServer(httpsPort)
	go httpsServer.ListenAndServeTLS(pemPath, keyPath)
}

func newServer(port string) *http.Server {
	return &http.Server{
		Addr:    ":" + port,
		Handler: http.HandlerFunc(proxyHandler),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	//simple AntiRobot
	if r.Header.Get("Naiba") != "lifelonglearning" {
		http.Error(w, "Naive", http.StatusForbidden)
		return
	}
	//get user IP
	var ip string
	ipSlice := strings.Split(r.RemoteAddr, ":")
	if len(ipSlice) == 2 {
		ip = ipSlice[0]
	} else {
		ip = "unknown"
	}
	//check domain limit
	if !proxyinabox.CheckIPDomain(ip, r.URL.Hostname()) {
		http.Error(w, "The request exceeds the limit, and up to "+strconv.Itoa(proxyinabox.DomainsPerIPHalfAnHour)+" domain names are crawled every half hour per IP.["+ip+"]", http.StatusForbidden)
		return
	}
	//set response header
	w.Header().Add("X-Powered-By", "Naiba")
	//dispath http request
	if r.Method == http.MethodConnect {
		handleTunneling("localhost:1087", w, r)
	} else {
		handleHTTP("http://localhost:1087", w, r)
	}
}
