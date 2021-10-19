package scraper

import (
	"better-av-tool/pkg/config"
)

type args struct {
	query string
	url   string
}

type testCase struct {
	name    string
	args    args
	wantErr bool
	want    string
}

func BeforeTest() {
	Setup(config.Proxy{
		Enable: true,
		Socket: "socks5://192.168.0.110:7891",
	})
}
