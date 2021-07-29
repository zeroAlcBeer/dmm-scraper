package scraper

import (
	myclient "better-av-tool/pkg/client"
	"better-av-tool/pkg/logger"
	"github.com/PuerkitoBio/goquery"
	"testing"
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

func Setup() {
	fc2Tests = make(map[string]*Fc2Scraper)
	mgsTests = make(map[string]*MGStageScraper)
	heyzoTests = make(map[string]*HeyzoScraper)
	gyuttoTests = make(map[string]*GyuttoScraper)
	dmmTests = make(map[string]*DMMScraper)

	log = logger.New()

	client = myclient.New()
	client.SetProxyUrl("socks5://192.168.0.110:7891")

	Init(client, log)
}

func TestDMMScraper_FetchDoc(t *testing.T) {
	tests := []testCase{
		{
			name: "fetchDoc expects no error",
			args: args{
				query: "gne-218",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, ok := dmmTests[tt.args.query]; !ok {
				dmmTests[tt.args.query] = &DMMScraper{}
			}
			if err := dmmTests[tt.args.query].FetchDoc(tt.args.query); (err != nil) != tt.wantErr {
				t.Errorf("FetchDoc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDMMScraper_GetActors(t *testing.T) {
	tests := []testCase{
		{
			name: "gne-218 expects actors",
			args: args{
				query: "gne-218",
			},
			want: "鈴村あいり",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dmmTests[tt.args.query].GetActors(); len(got) > 0 && got[0] != tt.want {
				t.Errorf("GetActors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDMMScraper_GetCover(t *testing.T) {
	tests := []testCase{
		{
			name: "gne-218 expects cover",
			args: args{
				query: "gne-218",
			},
			want: "https://pics.dmm.co.jp/mono/movie/adult/h_479gne218/h_479gne218pl.jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dmmTests[tt.args.query].GetCover(); got != tt.want {
				t.Errorf("GetCover() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDMMScraper_GetDirector(t *testing.T) {
	tests := []testCase{
		{
			name: "gne-218 expects director",
			args: args{
				query: "gne-218",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dmmTests[tt.args.query].GetDirector(); got != tt.want {
				t.Errorf("GetDirector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDMMScraper_GetLabel(t *testing.T) {
	tests := []testCase{
		{
			name: "gne-218 expects label",
			args: args{
				query: "gne-218",
			},
			want: "NEO GIFT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dmmTests[tt.args.query].GetLabel(); got != tt.want {
				t.Errorf("GetLabel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDMMScraper_GetMaker(t *testing.T) {
	tests := []testCase{
		{
			name: "gne-218 expects maker",
			args: args{
				query: "gne-218",
			},
			want: "GALLOP",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dmmTests[tt.args.query].GetMaker(); got != tt.want {
				t.Errorf("GetMaker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDMMScraper_GetNumber(t *testing.T) {
	tests := []testCase{
		{
			name: "gne-218 expects number",
			args: args{
				query: "gne-218",
			},
			want: "h_479gne218",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dmmTests[tt.args.query].GetNumber(); got != tt.want {
				t.Errorf("GetNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDMMScraper_GetPlot(t *testing.T) {
	tests := []testCase{
		{
			name: "gne-218 expects plot",
			args: args{
				query: "gne-218",
			},
			want: "「NEO GIFT」第ニ百十八弾！温泉旅行へ美少女たちと一泊二日でHなハメ撮りデート！開放感あふれる秘湯に浸かりながらハメハメしたり、絶品料理に舌鼓したあともホロ酔い気分でもう一発と美少女たちがいやらしく乱れる姿をご堪能下さい。",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dmmTests[tt.args.query].GetPlot(); got != tt.want {
				t.Errorf("GetPlot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDMMScraper_GetPremiered(t *testing.T) {
	tests := []testCase{
		{
			name: "gne-218 expects premiered",
			args: args{
				query: "gne-218",
			},
			want: "2019-01-04",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRel := dmmTests[tt.args.query].GetPremiered(); gotRel != tt.want {
				t.Errorf("GetPremiered() = %v, want %v", gotRel, tt.want)
			}
		})
	}
}

func TestDMMScraper_GetRuntime(t *testing.T) {
	tests := []testCase{
		{
			name: "gne-218 expects runtime",
			args: args{
				query: "gne-218",
			},
			want: "240分",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dmmTests[tt.args.query].GetRuntime(); got != tt.want {
				t.Errorf("GetRuntime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDMMScraper_GetSeries(t *testing.T) {
	tests := []testCase{
		{
			name: "gne-218 expects series",
			args: args{
				query: "gne-218",
			},
			want: "新・美少女貸切温泉旅行",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dmmTests[tt.args.query].GetSeries(); got != tt.want {
				t.Errorf("GetSeries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDMMScraper_GetTags(t *testing.T) {
	tests := []testCase{
		{
			name: "gne-218 expects tags",
			args: args{
				query: "gne-218",
			},
			want: "美少女",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dmmTests[tt.args.query].GetTags(); len(got) > 0 && got[0] != tt.want {
				t.Errorf("GetTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDMMScraper_GetTitle(t *testing.T) {
	tests := []testCase{
		{
			name: "gne-218 expects title",
			args: args{
				query: "gne-218",
			},
			want: "新・美少女貸切温泉旅行 5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dmmTests[tt.args.query].GetTitle(); got != tt.want {
				t.Errorf("GetTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
