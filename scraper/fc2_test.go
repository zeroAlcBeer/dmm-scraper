package scraper

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var fc2Tests map[string]*Fc2Scraper

type fields struct {
	doc        *goquery.Document
	docUrl     string
	HTTPClient *http.Client
}

type args struct {
	num string
}

func TestFc2Scraper_FetchDoc(t *testing.T) {
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
		{
			name: "fetchDoc expects no error",
			args: args{
				num: "559226",
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
					HTTPClient: proxyClient,
				}
			}
			if err := fc2Tests[tt.args.num].FetchDoc(tt.args.num); (err != nil) != tt.wantErr {
				t.Errorf("Fc2Scraper.FetchDoc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFc2Scraper_GetTitle(t *testing.T) {
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
		{
			name: "forceFetch expects no error",
			args: args{
				num: "559226",
			},
			want: "１８歳Jカップグラドル超人気美爆乳美女再度降臨。ハプニングありの期間枚数限定。後編",
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
		{
			name: "forceFetch expects no error",
			args: args{
				num: "559226",
			},
			want: "素人好きな親父ナンパ師",
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
