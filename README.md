# Proxy-in-a-Box
[![Go Report Card](https://goreportcard.com/badge/github.com/naiba/proxyinabox)](https://goreportcard.com/report/github.com/naiba/proxyinabox)[![travis](https://travis-ci.com/naiba/proxyinabox.svg?branch=master)](https://travis-ci.com/naiba/proxyinabox)

Proxy-in-a-Box helps programmers quickly and easily develop powerful crawler services. one-script, easy-to-use: proxies in a box.

# Usage
HTTP proxy: `http://localhost:8080`

HTTPS proxy: `https://localhost:8081`

Set in the code, and then grab it, the **Proxy-in-a-Box** will automatically assign the proxy.

# Preview
```
~$ cat test.sh
for ((i=1; i<=100; i ++))
do
    echo `curl --proxy https://localhost:8081 --proxy-insecure --proxy-header "Naiba: lifelonglearning" http://api.ip.la`
done
~$ ./test.sh
185.8.151.142
92.247.142.14
217.114.111.34
201.221.128.27
read tcp 10.5.1.187:64639->5.133.24.161:8081: read: connection reset by peer
85.21.240.153
104.237.227.68
89.218.223.250
```