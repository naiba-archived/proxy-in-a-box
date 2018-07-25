package main

import (
	"crypto/tls"

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

		IsDirect: true,
	}
	m.Init()
	m.ServeHTTP()
	select {}
}
