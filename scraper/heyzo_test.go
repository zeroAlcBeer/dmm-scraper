package scraper

import (
	"testing"
)
var heyzoTests map[string]*HeyzoScraper

func TestHeyzoScraper_FetchDoc(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "fetchDoc expects no error",
			args: args{
				num: "heyzo-1031",
			},
			wantErr: false,
		},
		{
			name: "fetchDoc expects no error",
			args: args{
				num: "heyzo-2169",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, ok := heyzoTests[tt.args.num]; !ok {
				heyzoTests[tt.args.num] = &HeyzoScraper{
					doc:        tt.fields.doc,
					docUrl:     tt.fields.docUrl,
					HTTPClient: proxyClient,
				}
			}
			if err := heyzoTests[tt.args.num].FetchDoc(tt.args.num); (err != nil) != tt.wantErr {
				t.Errorf("HeyzoScraper.FetchDoc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHeyzoScraper_GetTitle(t *testing.T) {
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1031 expects title",
			args: args{
				num: "heyzo-1031",
			},
			want: "試着がてらに試SEX！",
		},
		{
			name: "2169 expects title",
			args: args{
				num: "heyzo-2169",
			},
			want: "飲み過ぎ女たちとズッコンバッコン！",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := heyzoTests[tt.args.num].GetTitle(); got != tt.want {
				t.Errorf("HeyzoScraper.GetTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
