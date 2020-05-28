package scraper

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/proxy"

	"better-av-tool/log"
)

var dmmTests map[string]*DMMScraper
var fc2Tests map[string]*Fc2Scraper
var mgsTests map[string]*MGStageScraper
var heyzoTests map[string]*HeyzoScraper
var gyuttoTests map[string]*GyuttoScraper

type fields struct {
	doc *goquery.Document
}

type args struct {
	query string
	url   string
}

type testCase struct {
	name    string
	fields  fields
	args    args
	wantErr bool
	want    string
}

func TestMain(m *testing.M) {
	fc2Tests = make(map[string]*Fc2Scraper)
	mgsTests = make(map[string]*MGStageScraper)
	heyzoTests = make(map[string]*HeyzoScraper)
	gyuttoTests = make(map[string]*GyuttoScraper)
	dmmTests = make(map[string]*DMMScraper)

	u, _ := url.Parse("socks5://127.0.0.1:7891")
	dialer, err := proxy.FromURL(u, proxy.Direct)
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
