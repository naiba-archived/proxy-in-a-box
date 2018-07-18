# Proxy-in-a-Box
[![Go Report Card](https://goreportcard.com/badge/github.com/naiba/proxyinabox)](https://goreportcard.com/report/github.com/naiba/proxyinabox) [![travis](https://travis-ci.com/naiba/proxyinabox.svg?branch=master)](https://travis-ci.com/naiba/proxyinabox)

Proxy-in-a-Box helps programmers quickly and easily develop powerful crawler services. one-script, easy-to-use: proxies in a box.
```shell
Proxy-in-a-Box helps programmers quickly and easily develop powerful crawler services. one-script, easy-to-use: proxies in a box.

Usage:
  proxy-in-a-box [flags]

Flags:
  -c, --conf string   config file (default "./pb.yaml")
  -h, --help          help for proxy-in-a-box
  -p, --hp string     http proxy server port (default "8080")
  -s, --sp string     https proxy server port (default "8081")
```

## Usage
1. get lastest Proxy-in-a-Box
    ```
    go get -u -v github.com/naiba/proxyinabox/cmd/proxy-in-a-box/...
    ```
2. enter the application directory
    ```
    cd $GOPATH/bin
    ```
3. write config file #Config
4. run it
    ```
    ./proxy-in-a-box
    ```
5. configured in your code
    ```
    HTTP proxy: `http://[IP]:8080`
    HTTPS proxy: `https://[IP]:8081`
    * Please set http header when requesting: "Naiba: lifelonglearning" ref:https://github.com/naiba/proxyinabox/blob/master/cmd/proxy-in-a-box/test_server.sh
    ```
    Set in the code, and then grab it, the **Proxy-in-a-Box** will automatically assign the proxy.

## Config
```yaml
# run in debug mode
debug: false
# database config
db:
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
  # how many request can do per ip in 1min
  request_limit_per_ip: 420
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
    Latency     2.25s     3.02s   20.56s    90.69%
    Req/Sec     0.49      0.66     3.00     94.58%
  203 requests in 1.00m, 3.62MB read
  Non-2xx or 3xx responses: 7
Requests/sec:      3.38
Transfer/sec:     61.61KB
~$ wrk -H "Naiba: lifelonglearning"  -t30 -c30 -d60s -s proxy.lua --timeout 30s http://127.0.0.1:8080
Running 1m test @ http://127.0.0.1:8080
  30 threads and 30 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     5.45s     5.81s   24.13s    82.58%
    Req/Sec     0.96      1.85    10.00     77.78%
  162 requests in 1.00m, 55.74KB read
  Socket errors: connect 0, read 10, write 0, timeout 8
  Non-2xx or 3xx responses: 36
Requests/sec:      2.70
Transfer/sec:      0.93KB
```