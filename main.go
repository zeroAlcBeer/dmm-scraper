package main

import (
	"better-av-tool/scraper"
	"context"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"net/http"
)

var proxyClient *http.Client

func init() {
	proxyClient = &http.Client{}
	dialer, _ := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
	proxyClient.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			c, e := dialer.Dial(network, addr)
			return c, e
		},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func main() {
	s := &scraper.DMMScraper{}
	s.HTTPClient = proxyClient
	err := s.FetchDoc("saba-436")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s.GetPlot())
	fmt.Println(s.GetTitle())
	fmt.Println(s.GetDirector())
	fmt.Println(s.GetRuntime())
	fmt.Println(s.GetTags())
	fmt.Println(s.GetStudio())
	fmt.Println(s.GetMaker())
	fmt.Println(s.GetOutline())
	fmt.Println(s.GetActors())
	fmt.Println(s.GetLabel())
	fmt.Println(s.GetNumber())
	fmt.Println(s.GetCover())
	fmt.Println(s.GetWebsite())
	fmt.Println(s.GetPremiered())
	fmt.Println(s.GetSeries())

}
