package scraper

import "testing"

func TestGyuttoScraper_FetchDoc(t *testing.T) {
	BeforeTest()
	tests := []testCase{
		{
			name: "fetchDoc expects no error",
			args: args{
				query: "item240458",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GyuttoScraper{}
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
		})
	}
}