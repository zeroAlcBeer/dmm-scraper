package scraper

import (
	"testing"
)

func TestFc2Scraper_FetchDoc(t *testing.T) {
	tests := []testCase{
		{
			name: "fetchDoc expects no error",
			args: args{
				query: "1027251",
			},
			wantErr: false,
		},
		{
			name: "fetchDoc expects no error",
			args: args{
				query: "559226",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, ok := fc2Tests[tt.args.query]; !ok {
				fc2Tests[tt.args.query] = &Fc2Scraper{}
			}
			if err := fc2Tests[tt.args.query].FetchDoc(tt.args.query); (err != nil) != tt.wantErr {
				t.Errorf("Fc2Scraper.FetchDoc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFc2Scraper_GetTitle(t *testing.T) {
	tests := []testCase{
		{
			name: "1027251 expects title",
			args: args{
				query: "1027251",
			},
			want: "【某大手受付嬢】超絶イイ女！田中み○実似！Fカップありさちゃん（22才）と仕事終わりにそのままスーツでエッチww",
		},
		{
			name: "559226 expects title",
			args: args{
				query: "559226",
			},
			want: "１８歳Jカップグラドル超人気美爆乳美女再度降臨。ハプニングありの期間枚数限定。後編",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fc2Tests[tt.args.query].GetTitle(); got != tt.want {
				t.Errorf("Fc2Scraper.GetTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFc2Scraper_GetDirector(t *testing.T) {
	tests := []testCase{
		{
			name: "1027251 expects director",
			args: args{
				query: "1027251",
			},
			want: "ビッチとオレたち",
		},
		{
			name: "559226 expects director",
			args: args{
				query: "559226",
			},
			want: "素人好きな親父ナンパ師",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fc2Tests[tt.args.query].GetDirector(); got != tt.want {
				t.Errorf("Fc2Scraper.GetDirector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFc2Scraper_GetTags(t *testing.T) {
	tests := []testCase{
		{
			name: "1027251 expects tags",
			args: args{
				query: "1027251",
			},
			want: "ハメ撮り",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fc2Tests[tt.args.query].GetTags(); len(got) > 0 && got[0] != tt.want {
				t.Errorf("Fc2Scraper.GetTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFc2Scraper_GetNumber(t *testing.T) {
	tests := []testCase{
		{
			name: "1027251 expects number",
			args: args{
				query: "1027251",
			},
			want: "1027251",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fc2Tests[tt.args.query].GetNumber(); got != tt.want {
				t.Errorf("Fc2Scraper.GetNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFc2Scraper_GetCover(t *testing.T) {
	tests := []testCase{
		{
			name: "1027251 expects cover",
			args: args{
				query: "1027251",
			},
			want: "http://storage11000.contents.fc2.com/file/353/35243168/1548959337.39.jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fc2Tests[tt.args.query].GetCover(); got != tt.want {
				t.Errorf("Fc2Scraper.GetCover() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFc2Scraper_GetPremiered(t *testing.T) {
	tests := []testCase{
		{
			name: "1027251 expects date",
			args: args{
				query: "1027251",
			},
			want: "2019-02-02",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fc2Tests[tt.args.query].GetPremiered(); got != tt.want {
				t.Errorf("Fc2Scraper.GetPremiered() = %v, want %v", got, tt.want)
			}
		})
	}
}
