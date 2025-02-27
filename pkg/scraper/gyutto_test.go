package scraper

import "testing"

func TestGyuttoScraper_FetchDoc(t *testing.T) {
	BeforeTest()
	tests := []testCase{
		{
			name: "fetchDoc expects no error",
			args: args{
				query: "item207303",
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
			dumpS(t, s)
		})
	}
}
