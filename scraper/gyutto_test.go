package scraper

import (
	"testing"
)

func TestGyuttoScraper_FetchDoc(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "fetchDoc expects no error",
			args: args{
				query: "FG0礼装マシュちゃん[HF]",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, ok := gyuttoTests[tt.args.query]; !ok {
				gyuttoTests[tt.args.query] = &GyuttoScraper{}
			}
			if err := gyuttoTests[tt.args.query].FetchDoc(tt.args.query, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("FetchDoc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGyuttoScraper_GetTitle(t *testing.T) {
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "getTitle expects no error",
			args: args{
				query: "FG0礼装マシュちゃん[HF]",
			},
			want: "ガチ本壊ちゃん　過去最高感度の過去最高おっぱい 10代J○から調教済みオフパコびっちレイヤー　FG0礼装マシュちゃん[HF]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gyuttoTests[tt.args.query].GetTitle(); got != tt.want {
				t.Errorf("GetTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
