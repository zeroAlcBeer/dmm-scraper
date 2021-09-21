package scraper

import (
	"better-av-tool/pkg/config"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"testing"
)

func TestDefaultScraper_GetDocFromURL(t *testing.T) {
	Setup(config.Proxy{
		Enable: true,
		Socket: "socks5://192.168.0.110:7891",
	})
	type fields struct {
		doc       *goquery.Document
		cookie    *http.Cookie
		isArchive bool
	}
	type args struct {
		u string
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
				u: "https://gyutto.com/i/item240458?select_uaflag=1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GyuttoScraper{}
			if err := s.GetDocFromURL(tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("GetDocFromURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("%v", s.doc)
			t.Logf("%v", s.GetTitle())
		})
	}
}
