# Proxy-in-a-Box

[![Go Report Card](https://goreportcard.com/badge/github.com/naiba/proxyinabox)](https://goreportcard.com/report/github.com/naiba/proxyinabox) [![travis](https://travis-ci.com/naiba/proxyinabox.svg?branch=master)](https://travis-ci.com/naiba/proxyinabox)

Proxy-in-a-Box helps programmers quickly and easily develop powerful crawler services. one-script, easy-to-use: proxies in a box.

```shell
Usage:
  proxy-in-a-box [flags]

Flags:
  -c, --conf string   config file (default "./pb.yaml")
  -h, --help          help for proxy-in-a-box
  -p, --ha string     http proxy server addr (default "8080")
  -s, --sa string     https proxy server addr (default "8081")
```

## Usage

1. get lastest Proxy-in-a-Box
    ```shell
    go get -u -v github.com/naiba/proxyinabox/cmd/proxy-in-a-box/...
    ```
2. enter the application directory
    ```shell
    cd $GOPATH/bin
    ```
3. write config file #Config
4. run it
    ```shell
    ./proxy-in-a-box
    ```
5. configured in your code
    ```none
    HTTP proxy: `http://[IP]:8080`
    HTTPS proxy: `https://[IP]:8081`
    * Please set http header when requesting: "Naiba: lifelonglearning" ref:https://github.com/naiba/proxyinabox/blob/master/cmd/proxy-in-a-box/test_server.sh
    ```
    Set in the code, and then grab it, the **Proxy-in-a-Box** will automatically assign the proxy.

## Config

```yaml
# run in debug mode
debug: true
# mysql config
mysql:
  host: 127.0.0.1
  port: 3306
  user: root
  pass: 123456
  dbname: proxy
# system config
sys:
  name: Naiba
  # verify proxy's worker num
  proxy_verify_worker: 20
  # how many domains can request per ip in 30min
  domains_per_ip: 30
  # how many request can do per ip in 1s
  request_limit_per_ip: 10
  # verify interval of the proxy stored in the database
  verify_duration: 30
```

## Benchmark

Server bandwidth and mysql configuration will affect the test results, mysql configuration affects the scheduling of the agent.

```shell
~$ wrk -H "Naiba: lifelonglearning"  -t30 -c30 -d60s -s proxy.lua --timeout 30s http://127.0.0.1:8080
Running 1m test @ http://127.0.0.1:8080
  30 threads and 30 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.15s     2.76s   23.61s    88.89%
    Req/Sec    14.17      8.69    30.00     70.42%
  1058 requests in 1.00m, 487.21KB read
  Socket errors: connect 0, read 7, write 0, timeout 10
  Non-2xx or 3xx responses: 37
Requests/sec:     17.61
Transfer/sec:      8.11KB
~$ wrk -H "Naiba: lifelonglearning"  -t50 -c50 -d60s -s proxy.lua --timeout 30s http://127.0.0.1:8080
Running 1m test @ http://127.0.0.1:8080
  50 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.50s     3.44s   24.30s    88.89%
    Req/Sec    13.01      9.95    30.00     58.77%
  1050 requests in 1.00m, 500.08KB read
  Socket errors: connect 0, read 15, write 0, timeout 15
  Non-2xx or 3xx responses: 50
Requests/sec:     17.47
Transfer/sec:      8.32KB
```