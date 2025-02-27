package scraper

import (
	"testing"
)

func TestAVWikiScraper_FetchDoc(t *testing.T) {
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
			s := &AVWikiScraper{}
			if err := s.FetchDoc(tt.args.query); (err != nil) != tt.wantErr {
				t.Errorf("FetchDoc() error = %v, wantErr %v", err, tt.wantErr)
			}
			dumpS(t, s)
		})
	}
}
