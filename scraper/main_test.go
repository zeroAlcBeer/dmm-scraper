package scraper

import (
	"better-av-tool/log"
	"context"
	"crypto/tls"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"testing"
)

var proxyClient *http.Client

func TestMain(m *testing.M) {
	fc2Tests = make(map[string]*Fc2Scraper)
	mgsTests = make(map[string]*MGStageScraper)
	heyzoTests = make(map[string]*HeyzoScraper)

	url, err := url.Parse("socks5://127.0.0.1:7891")
	if err != nil {
		log.Fatal(err)
	}
	dialer, err := proxy.FromURL(url, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}
	proxyClient = &http.Client{}
	proxyClient.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			c, e := dialer.Dial(network, addr)
			return c, e
		},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	m.Run()
}
