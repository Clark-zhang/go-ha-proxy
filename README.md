# Introduction:
  This go-ha-proxy is simple reverse proxy using Golang.

  Original fork from [matishsiao/go-ha-proxy](https://github.com/matishsiao/go-ha-proxy)

  I have reorganize the project from matishsia to make it work on go1.9.2. Also add some test to it. You could check **How to test** section to see how to test this project manually. Hope this project will help you to learn golang as it did for me. - Clark Zhang

## Version:

version:0.0.2

1.support TCP proxy and you can use round-robin or source mode to HA.

2.health mode supported.

3.monitor server supported.

    http://localip:8080

4.command line params supported.

    ./GoHAProxy -help

5.auto reload config file.When file changed.GoHAProxy will auto reloaded.

6.balance type:RoundRobin,Source,Weight,LeastConn

    RoundRobin:auto RR

    Source:use source ip to balance service.

    Weight:Use weight setting to balance.

    LeastConn:Search least conn server to balance server.


## Install:
```sh
  go get github.com/Clark-zhang/go-ha-proxy
```

## How to test (by Clark Zhang)

### start two server
cd test

go run go_http_router.go -port=9001

go run go_http_router.go -port=9002

### start ha-proxy
go run goHAProxy.go -debug=true

### try to fetch data or kill one of the server
curl localhost:9000


## Configuration format:

configs in config.json

ProxyList:

      Proxy:

      Src:source ip or domain.

      SrcPort:source port.

      Mode:tcp,health

      Type:RoundRobin,Source,Weight,LeastConn

      KeepAlive:1 second (keep alive server connection.Set zero will allways keep alive.)

      CheckTime:1 second (check health,default value:5 seconds.)

      DstList:Destination server list

         DstNode:

         Name:server name

         Dst:server ip or domain

         DstPort:destination port

         Weight:not use(when Weight HA mode done,will use this arg)

         Check:true or false(check node server health.If you set false,the GoHAProxy will set this server allways health.)
