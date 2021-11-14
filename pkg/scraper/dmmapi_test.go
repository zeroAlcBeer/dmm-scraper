package scraper

import (
	"testing"
)

func TestDMMApiScraper_FetchDoc(t *testing.T) {
	BeforeTest()
	tests := []testCase{
		{
			name: "fetchDoc expects no error",
			args: args{
				query: "ACRN-119",
			},
			wantErr: false,
		},
		{
			name: "fetchDoc expects no error",
			args: args{
				query: "pred-352",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DMMApiDigitalScraper{}
			if err := s.FetchDoc(tt.args.query); (err != nil) != tt.wantErr {
				t.Errorf("FetchDoc() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := s.GetNumber()
			t.Logf("GetNumber() = %v", got)
			got = s.GetPlot()
			t.Logf("GetPlot() = %v", got)
			got = s.GetTitle()
			t.Logf("GetTitle() = %v", got)
			got = s.GetCover()
			t.Logf("GetCover() = %v", got)
			gots := s.GetActors()
			t.Logf("GetActors() = %v", gots)
		})
	}
}
