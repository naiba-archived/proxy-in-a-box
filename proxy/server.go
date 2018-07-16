package proxy

import (
	"fmt"
	"net/http"

	"github.com/elazarl/goproxy"
)

//Serv serv proxy server
func Serv(port string) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	fmt.Println(http.ListenAndServe(":"+port, proxy))
	proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return r, goproxy.NewResponse(r,
			goproxy.ContentTypeText, http.StatusForbidden,
			"Don't waste your time!")
	})
}
