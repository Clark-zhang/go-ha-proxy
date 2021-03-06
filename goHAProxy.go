package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	haProxy "github.com/Clark-zhang/go-ha-proxy/src/ha-proxy"
)

var version string = "0.0.2"
var proxyServer haProxy.ProxyServerStruct
var configInfo haProxy.ConfigInfoStruct
var help *bool

func main() {
	configInfo.FileName = flag.String("config", "./config.json", "set config file path.")
	configInfo.Debug = flag.Bool("debug", false, "show debug trace message.")
	configInfo.Version = flag.Bool("version", false, "GoHAProxy version.")
	help = flag.Bool("help", false, "Show help information.")
	flag.Parse()

	if *help {
		fmt.Printf("GoHAProxy is simple HAProxy by Golang,You can use this to load balance TCP or check health status.\n")
		fmt.Println("GoHAProxy Version", version)
		fmt.Printf("Monitor server:local ip:8080, Default value:[:8080]\n")
		fmt.Printf("-config    Set cofing file path. Default value:%v\n", *configInfo.FileName)
		fmt.Printf("-debug     Show debug trace message. Default value:%v\n", *configInfo.Debug)
		fmt.Printf("-version   Show GoHAProxy version.\n")
		fmt.Printf("-help      Show help information.\n")
		os.Exit(0)
	}

	if *configInfo.Version {
		fmt.Println("GoHAProxy Version", version)
		os.Exit(0)
	}
	ok,haConfig := loadConfigs()
	if ok {
		for _, proxy := range haConfig.Configs.ProxyList {
			FS := new(haProxy.ForwardServer)
			proxyServer.ServerList = append(proxyServer.ServerList, FS)

			//依赖注入
			haProxy.ProxyServer = proxyServer

			go FS.Listen(proxy)
		}
	}

	//依赖注入
	haProxy.ConfigInfo = configInfo

	go haProxy.Monitor()
	for {
		configWatcher()
		time.Sleep(500 * time.Millisecond)
	}
}

func configWatcher() {
	file, err := os.Open(*configInfo.FileName) // For read access.
	if err != nil {
		fmt.Println(err)
	}
	info, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}
	if configInfo.Size == 0 {
		configInfo.Size = info.Size()
		configInfo.ModTime = info.ModTime()
	}

	if info.Size() != configInfo.Size || info.ModTime() != configInfo.ModTime {
		fmt.Printf("Config changed.Reolad.\n")
		configInfo.Size = info.Size()
		configInfo.ModTime = info.ModTime()
		ok,haConfig := loadConfigs()
		if ok {
			//檢查有沒有移除掉的設定有的話就移除
			for k, sProxy := range proxyServer.ServerList {
				oldProxy := true
				for _, proxy := range haConfig.Configs.ProxyList {
					if *configInfo.Debug {
						fmt.Printf("Delete PName:%s SrvProxyName:%s \n", proxy.Name, sProxy.SrvProxy.Name)
					}
					if proxy.Name == sProxy.SrvProxy.Name {
						oldProxy = false
						break
					}
				}

				if oldProxy {
					if *configInfo.Debug {
						fmt.Printf("Delete Proxy[%v]: %v\n", k, sProxy.SrvProxy.Name)
					}
					sProxy.Stop()
					proxyServer.ServerList = proxyServer.ServerList[:k+copy(proxyServer.ServerList[k:], proxyServer.ServerList[k+1:])]

				}
			}

			//檢查有沒有新的設定
			for _, proxy := range haConfig.Configs.ProxyList {
				newProxy := true
				for k, sProxy := range proxyServer.ServerList {
					if *configInfo.Debug {
						fmt.Printf("Check Add PName[%d]:%s SrvProxyName:%s \n", k, proxy.Name, sProxy.SrvProxy.Name)
					}
					if proxy.Name == sProxy.SrvProxy.Name {
						proxyServer.ServerList[k].Reload(proxy)
						newProxy = false
						break
					}
				}
				if newProxy {
					FS := new(haProxy.ForwardServer)
					proxyServer.ServerList = append(proxyServer.ServerList, FS)
					if *configInfo.Debug {
						fmt.Printf("Add New Proxy: %v\n", proxy)
					}
					go FS.Listen(proxy)
				}
			}
		}

	}
	defer file.Close()
}

func loadConfigs() (bool, haProxy.HAConfig) {
	fmt.Printf("GoHAProxy load config file:%s\n", *configInfo.FileName)
	file, e := ioutil.ReadFile(*configInfo.FileName)
	if e != nil {
		fmt.Printf("Load GoHAProxy config error: %v\n", e)
		os.Exit(1)
	}
	var haConfig haProxy.HAConfig
	err := json.Unmarshal(file, &haConfig)
	if err != nil {
		fmt.Printf("Config load error:%v \n",err)
		return false,haConfig
	}
	return true,haConfig
}
