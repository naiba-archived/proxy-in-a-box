package proxy

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/naiba/proxyinabox"
	"github.com/naiba/proxyinabox/service/sqlite3"
)

var domainService proxyinabox.DomainService
var activityService proxyinabox.ActivityService

//Serv serv the http proxy
func Serv(httpPort, httpsPort string) {
	//init service
	domainService = &sqlite3.DomainService{DB: proxyinabox.DB}
	activityService = &sqlite3.ActivityService{DB: proxyinabox.DB}

	//start http proxy server
	httpServer := newServer(httpPort)
	go httpServer.ListenAndServe()

	//start https proxy server
	var pemPath = "./server.pem"
	var keyPath = "./server.key"
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
	var domain = r.URL.Hostname()
	if !proxyinabox.CheckIPDomain(ip, domain) {
		http.Error(w, "The request exceeds the limit, and up to "+strconv.Itoa(proxyinabox.DomainsPerIPHalfAnHour)+" domain names are crawled every half hour per IP.["+ip+"]", http.StatusForbidden)
		return
	}
	//get a proxy
	p, f := getProxy(r.Method, domain)
	//set response header
	w.Header().Add("X-Powered-By", "Naiba")
	//dispath http request
	f(p, w, r)
	if r.Method == http.MethodConnect {
		handleTunneling("localhost:1087", w, r)
	} else {
		handleHTTP("http://localhost:1087", w, r)
	}
}

func getProxy(method, domain string) (string, func(p string, w http.ResponseWriter, r *http.Request)) {
	//get domain by name
	d, err := domainService.GetByName(domain)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			//unknown error
			fmt.Println("panic", err)
		} else {
			//new domain save to database
			proxyinabox.DB.Save(&proxyinabox.Domain{
				Name: domain,
			})
			//get a fresh proxy

		}
	}
}
