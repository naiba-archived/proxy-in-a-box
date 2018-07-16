package proxy

import (
	"crypto/tls"
	"fmt"
	"net/http"
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
		Addr: ":" + port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(port, r.Method, r.Host)
			w.Header().Add("X-Powered-By", "Naiba")
			if r.Method == http.MethodConnect {
				handleTunneling("localhost:1087", w, r)
			} else {
				handleHTTP("http://localhost:1087", w, r)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
}
