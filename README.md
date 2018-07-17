# Proxy-in-a-Box
[![Go Report Card](https://goreportcard.com/badge/github.com/naiba/proxyinabox)](https://goreportcard.com/report/github.com/naiba/proxyinabox) [![travis](https://travis-ci.com/naiba/proxyinabox.svg?branch=master)](https://travis-ci.com/naiba/proxyinabox)

Proxy-in-a-Box helps programmers quickly and easily develop powerful crawler services. one-script, easy-to-use: proxies in a box.

# Usage
1. get Proxy-in-a-Box
    ```
    go get -v github.com/naiba/proxyinabox/cmd/proxy-in-a-box/...
    ```
2. run Proxy-in-a-Box
    ```
    cd $GOPATH/bin
    ./proxy-in-a-box
    ```
3. setup your application
    ```
    HTTP proxy: `http://localhost:8080`
    HTTPS proxy: `https://localhost:8081`
    ```
    Set in the code, and then grab it, the **Proxy-in-a-Box** will automatically assign the proxy.

# Demo
Contact owner.