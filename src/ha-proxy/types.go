package haProxy

import (
	"time"
)

type ConfigInfoStruct struct {
	FileName *string
	Size     int64
	ModTime  time.Time
	Debug    *bool
	Version  *bool
}

type ProxyServerStruct struct {
	ServerList []*ForwardServer
}

type HtmlData struct {
	Title       string
	Data        map[string]string
	ProxyStatus []Proxy
}

type HAConfig struct {
	Configs Config
}

type Config struct {
	ProxyList []Proxy
}

type Proxy struct {
	Name      string
	Src       string
	SrcPort   string
	Mode      string
	Type      string
	KeepAlive int
	CheckTime int
	Counter   int
	DstLen    int
	Index     int
	DstList   []DstConfig
}

type DstConfig struct {
	Name        string
	Dst         string
	DstPort     string
	Weight      int
	WeightCounter      int
	Check       bool
	Health      bool
	Counter     int
	Connections int
}

var ConfigInfo ConfigInfoStruct
var ProxyServer ProxyServerStruct

func (p *Proxy) GetSrcAddr() string {
	return p.Src + ":" + p.SrcPort
}

func (p *Proxy) GetDstAddr(index int) string {
	return p.DstList[index].Dst + ":" + p.DstList[index].DstPort
}
