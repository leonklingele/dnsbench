# DNSBench — DNS benchmarking utility

## Installation

```sh
go get -u github.com/leonklingele/dnsbench/...
dnsbench -help
```

## Run benchmark

```sh
dnsbench -domains google.com,cloudflare.com -queries 8 -servers 1.1.1.1,8.8.8.8

# Example output:
Domains: google.com, cloudflare.com
Servers: 1.1.1.1:53, 8.8.8.8:53
Proto:   udp
Queries: 8
Workers: 1
[1.1.1.1:53]: avg query time for google.com     : 19.945912ms
[1.1.1.1:53]: avg query time for cloudflare.com : 18.016792ms
[8.8.8.8:53]: avg query time for google.com     : 10.637054ms
[8.8.8.8:53]: avg query time for cloudflare.com : 10.906691ms
Summary [8.8.8.8:53]: avg query time: 10.771872ms
Summary [1.1.1.1:53]: avg query time: 18.981352ms
```
