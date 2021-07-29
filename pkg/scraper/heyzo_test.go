package scraper

import (
	"testing"
)

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
				query: "heyzo-1031",
			},
			wantErr: false,
		},
		{
			name: "fetchDoc expects no error",
			args: args{
				query: "heyzo-2169",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, ok := heyzoTests[tt.args.query]; !ok {
				heyzoTests[tt.args.query] = &HeyzoScraper{
					doc: tt.fields.doc,
				}
			}
			if err := heyzoTests[tt.args.query].FetchDoc(tt.args.query); (err != nil) != tt.wantErr {
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
				query: "heyzo-1031",
			},
			want: "試着がてらに試SEX！",
		},
		{
			name: "2169 expects title",
			args: args{
				query: "heyzo-2169",
			},
			want: "飲み過ぎ女たちとズッコンバッコン！",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := heyzoTests[tt.args.query].GetTitle(); got != tt.want {
				t.Errorf("HeyzoScraper.GetTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
