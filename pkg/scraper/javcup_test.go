package scraper

import (
	"testing"
)

func TestJavCupScraper_FetchDoc(t *testing.T) {
	BeforeTest()
	tests := []testCase{
		{
			name: "fetchDoc expects no error",
			args: args{
				query: "OYCVR020.VR",
			},
			wantErr: false,
			want:    "OYCVR-020",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &JavCupScraper{}
			if err := s.FetchDoc(tt.args.query); (err != nil) != tt.wantErr {
				t.Errorf("FetchDoc() error = %v, wantErr %v", err, tt.wantErr)
			}
			//debug, _ := s.doc.Html()
			//t.Logf(debug)
			got := s.GetNumber()
			t.Logf("GetNumber() = %v", got)
			got = s.GetPlot()
			t.Logf("GetPlot() = %v", got)
			got = s.GetTitle()
			t.Logf("GetTitle() = %v", got)
			got = s.GetCover()
			t.Logf("GetCover() = %v", got)
			got = s.GetDirector()
			t.Logf("GetDirector() = %v", got)
			got = s.GetMaker()
			t.Logf("GetMaker() = %v", got)
			gots := s.GetTags()
			t.Logf("GetTags() = %v", gots)
		})
	}
}
