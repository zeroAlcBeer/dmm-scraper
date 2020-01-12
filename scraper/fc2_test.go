package scraper

import (
	"better-av-tool/log"
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/proxy"
)

var fc2Tests map[string]*Fc2Scraper
var proxyClient *http.Client

func TestMain(m *testing.M) {
	fc2Tests = make(map[string]*Fc2Scraper)
	mgsTests = make(map[string]*MGStageScraper)

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

func TestFc2Scraper_FetchDoc(t *testing.T) {
	type fields struct {
		doc        *goquery.Document
		docUrl     string
		HTTPClient *http.Client
	}
	type args struct {
		num string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "fetchDoc expects no error",
			args: args{
				num: "1027251",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, ok := fc2Tests[tt.args.num]; !ok {
				fc2Tests[tt.args.num] = &Fc2Scraper{
					doc:        tt.fields.doc,
					docUrl:     tt.fields.docUrl,
					HTTPClient: tt.fields.HTTPClient,
				}
			}
			if err := fc2Tests[tt.args.num].FetchDoc(tt.args.num); (err != nil) != tt.wantErr {
				t.Errorf("Fc2Scraper.FetchDoc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFc2Scraper_GetTitle(t *testing.T) {
	type args struct {
		num string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1027251 expects title",
			args: args{
				num: "1027251",
			},
			want: "【某大手受付嬢】超絶イイ女！田中み○実似！Fカップありさちゃん（22才）と仕事終わりにそのままスーツでエッチww",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fc2Tests[tt.args.num].GetTitle(); got != tt.want {
				t.Errorf("Fc2Scraper.GetTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFc2Scraper_GetDirector(t *testing.T) {
	type args struct {
		num string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1027251 expects director",
			args: args{
				num: "1027251",
			},
			want: "ビッチとオレたち",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fc2Tests[tt.args.num].GetDirector(); got != tt.want {
				t.Errorf("Fc2Scraper.GetDirector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFc2Scraper_GetRuntime(t *testing.T) {
	type args struct {
		num string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1027251 expects runtime",
			args: args{
				num: "1027251",
			},
			want: "52:07",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fc2Tests[tt.args.num].GetRuntime(); got != tt.want {
				t.Errorf("Fc2Scraper.GetRuntime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFc2Scraper_GetTags(t *testing.T) {
	type args struct {
		num string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1027251 expects tags",
			args: args{
				num: "1027251",
			},
			want: "ハメ撮り",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fc2Tests[tt.args.num].GetTags(); got[0] != tt.want {
				t.Errorf("Fc2Scraper.GetTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFc2Scraper_GetNumber(t *testing.T) {
	type args struct {
		num string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1027251 expects number",
			args: args{
				num: "1027251",
			},
			want: "1027251",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fc2Tests[tt.args.num].GetNumber(); got != tt.want {
				t.Errorf("Fc2Scraper.GetNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFc2Scraper_GetCover(t *testing.T) {
	type args struct {
		num string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1027251 expects number",
			args: args{
				num: "1027251",
			},
			want: "https://storage11000.contents.fc2.com/file/353/35243168/1548959337.39.jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fc2Tests[tt.args.num].GetCover(); got != tt.want {
				t.Errorf("Fc2Scraper.GetCover() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFc2Scraper_GetPremiered(t *testing.T) {
	type args struct {
		num string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1027251 expects number",
			args: args{
				num: "1027251",
			},
			want: "2019/02/02",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fc2Tests[tt.args.num].GetPremiered(); got != tt.want {
				t.Errorf("Fc2Scraper.GetPremiered() = %v, want %v", got, tt.want)
			}
		})
	}
}
