package scraper

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
)
var mgsTests map[string]*MGStageScraper

func TestMGStageScraper_FetchDoc(t *testing.T) {
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
				num: "siro-3171",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, ok := mgsTests[tt.args.num]; !ok {
				mgsTests[tt.args.num] = &MGStageScraper{
					doc:        tt.fields.doc,
					docUrl:     tt.fields.docUrl,
					HTTPClient: proxyClient,
				}
			}
			if err := mgsTests[tt.args.num].FetchDoc(tt.args.num); (err != nil) != tt.wantErr {
				t.Errorf("FetchDoc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMGStageScraper_GetPlot(t *testing.T) {
	type args struct {
		num string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "siro-3171 expects title",
			args: args{
				num: "siro-3171",
			},
			want: "とても落ち着いた雰囲気のお姉さんタイプのあゆみさん。まるでモデルさんみたいです！普段はOLをしているらしいですが、こんな綺麗なお姉さんが仕事場にいるなんて、男はほっとかないでしょ！と思いきや全然誘われないそうです。逆に綺麗すぎてみんな敬遠してるのかなーと勝手に想像しちゃいますね。でも上司からはすれ違いざまにお尻などを触られるらしく、それは「相手が私と仲良くなりたいのかなって思っちゃいます」と女神のような返答。クールな美女タイプかと思いきや、中身はとても優しい暖かお姉さんなんですね。現在は彼氏はいないそうで、ムラムラしたときは専らオナニー！そんなオナニーを見せて下さい！って言ったら恥ずかしがりながらオナニーを披露してくれます。最初は緊張気味でしたが、いざ感じ始めるとその勢いは素晴らしいの一言！大胆に自分の秘部を弄りエロスイッチオン！目がトロンとして、男優のち○ぽを見つめながら熱い吐息を吐いておねだりしてきます。さっきまでは優しいお姉さんだったのに、スイッチが入ったらいきなりエロい妖艶なお姉さんに早変わり！このギャップがまたいいですね。色んな女性の表情が見えるって良いですね。",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mgsTests[tt.args.num].GetPlot(); got != tt.want {
				t.Errorf("GetPlot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMGStageScraper_GetTitle(t *testing.T) {
	type args struct {
		num string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "siro-3171 expects title",
			args: args{
				num: "siro-3171",
			},
			want: "【初撮り】ネットでAV応募→AV体験撮影 400",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mgsTests[tt.args.num].GetTitle(); got != tt.want {
				t.Errorf("GetTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestMGStageScraper_GetDirector(t *testing.T) {
//	type fields struct {
//		doc        *goquery.Document
//		docUrl     string
//		HTTPClient *http.Client
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &MGStageScraper{
//				doc:        tt.fields.doc,
//				docUrl:     tt.fields.docUrl,
//				HTTPClient: tt.fields.HTTPClient,
//			}
//			if got := s.GetDirector(); got != tt.want {
//				t.Errorf("MGStageScraper.GetDirector() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMGStageScraper_GetRuntime(t *testing.T) {
//	type fields struct {
//		doc        *goquery.Document
//		docUrl     string
//		HTTPClient *http.Client
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &MGStageScraper{
//				doc:        tt.fields.doc,
//				docUrl:     tt.fields.docUrl,
//				HTTPClient: tt.fields.HTTPClient,
//			}
//			if got := s.GetRuntime(); got != tt.want {
//				t.Errorf("MGStageScraper.GetRuntime() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMGStageScraper_GetTags(t *testing.T) {
//	type fields struct {
//		doc        *goquery.Document
//		docUrl     string
//		HTTPClient *http.Client
//	}
//	tests := []struct {
//		name     string
//		fields   fields
//		wantTags []string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &MGStageScraper{
//				doc:        tt.fields.doc,
//				docUrl:     tt.fields.docUrl,
//				HTTPClient: tt.fields.HTTPClient,
//			}
//			if gotTags := s.GetTags(); !reflect.DeepEqual(gotTags, tt.wantTags) {
//				t.Errorf("MGStageScraper.GetTags() = %v, want %v", gotTags, tt.wantTags)
//			}
//		})
//	}
//}
//
//func TestMGStageScraper_GetMaker(t *testing.T) {
//	type fields struct {
//		doc        *goquery.Document
//		docUrl     string
//		HTTPClient *http.Client
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &MGStageScraper{
//				doc:        tt.fields.doc,
//				docUrl:     tt.fields.docUrl,
//				HTTPClient: tt.fields.HTTPClient,
//			}
//			if got := s.GetMaker(); got != tt.want {
//				t.Errorf("MGStageScraper.GetMaker() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMGStageScraper_GetActors(t *testing.T) {
//	type fields struct {
//		doc        *goquery.Document
//		docUrl     string
//		HTTPClient *http.Client
//	}
//	tests := []struct {
//		name       string
//		fields     fields
//		wantActors []string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &MGStageScraper{
//				doc:        tt.fields.doc,
//				docUrl:     tt.fields.docUrl,
//				HTTPClient: tt.fields.HTTPClient,
//			}
//			if gotActors := s.GetActors(); !reflect.DeepEqual(gotActors, tt.wantActors) {
//				t.Errorf("MGStageScraper.GetActors() = %v, want %v", gotActors, tt.wantActors)
//			}
//		})
//	}
//}
//
//func TestMGStageScraper_GetLabel(t *testing.T) {
//	type fields struct {
//		doc        *goquery.Document
//		docUrl     string
//		HTTPClient *http.Client
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &MGStageScraper{
//				doc:        tt.fields.doc,
//				docUrl:     tt.fields.docUrl,
//				HTTPClient: tt.fields.HTTPClient,
//			}
//			if got := s.GetLabel(); got != tt.want {
//				t.Errorf("MGStageScraper.GetLabel() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMGStageScraper_GetNumber(t *testing.T) {
//	type fields struct {
//		doc        *goquery.Document
//		docUrl     string
//		HTTPClient *http.Client
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &MGStageScraper{
//				doc:        tt.fields.doc,
//				docUrl:     tt.fields.docUrl,
//				HTTPClient: tt.fields.HTTPClient,
//			}
//			if got := s.GetNumber(); got != tt.want {
//				t.Errorf("MGStageScraper.GetNumber() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMGStageScraper_GetCover(t *testing.T) {
//	type fields struct {
//		doc        *goquery.Document
//		docUrl     string
//		HTTPClient *http.Client
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &MGStageScraper{
//				doc:        tt.fields.doc,
//				docUrl:     tt.fields.docUrl,
//				HTTPClient: tt.fields.HTTPClient,
//			}
//			if got := s.GetCover(); got != tt.want {
//				t.Errorf("MGStageScraper.GetCover() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMGStageScraper_GetWebsite(t *testing.T) {
//	type fields struct {
//		doc        *goquery.Document
//		docUrl     string
//		HTTPClient *http.Client
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &MGStageScraper{
//				doc:        tt.fields.doc,
//				docUrl:     tt.fields.docUrl,
//				HTTPClient: tt.fields.HTTPClient,
//			}
//			if got := s.GetWebsite(); got != tt.want {
//				t.Errorf("MGStageScraper.GetWebsite() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMGStageScraper_GetPremiered(t *testing.T) {
//	type fields struct {
//		doc        *goquery.Document
//		docUrl     string
//		HTTPClient *http.Client
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &MGStageScraper{
//				doc:        tt.fields.doc,
//				docUrl:     tt.fields.docUrl,
//				HTTPClient: tt.fields.HTTPClient,
//			}
//			if got := s.GetPremiered(); got != tt.want {
//				t.Errorf("MGStageScraper.GetPremiered() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMGStageScraper_GetSeries(t *testing.T) {
//	type fields struct {
//		doc        *goquery.Document
//		docUrl     string
//		HTTPClient *http.Client
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &MGStageScraper{
//				doc:        tt.fields.doc,
//				docUrl:     tt.fields.docUrl,
//				HTTPClient: tt.fields.HTTPClient,
//			}
//			if got := s.GetSeries(); got != tt.want {
//				t.Errorf("MGStageScraper.GetSeries() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_getMgstageTableValue(t *testing.T) {
//	type args struct {
//		key string
//		doc *goquery.Document
//	}
//	tests := []struct {
//		name       string
//		args       args
//		wantTarget *goquery.Selection
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if gotTarget := getMgstageTableValue(tt.args.key, tt.args.doc); !reflect.DeepEqual(gotTarget, tt.wantTarget) {
//				t.Errorf("getMgstageTableValue() = %v, want %v", gotTarget, tt.wantTarget)
//			}
//		})
//	}
//}
