package main

import (
	"crypto/tls"
	"net/http"

	"github.com/naiba/proxyinabox/mitm"
)

func main() {
	m := &mitm.MITM{
		ListenHTTPS: true,
		HTTPAddr:    "127.0.0.1:8080",
		HTTPSAddr:   "127.0.0.1:8081",
		TLSConf: &struct {
			PrivateKeyFile  string
			CertFile        string
			Organization    string
			CommonName      string
			ServerTLSConfig *tls.Config
		}{
			PrivateKeyFile: "server.key",
			CertFile:       "server.pem",
		},

		IsDirect: false,
		Scheduler: func(r *http.Request) (string, error) {
			return "127.0.0.1:1087", nil
		},
	}
	m.Init()
	m.ServeHTTP()
	select {}
}
