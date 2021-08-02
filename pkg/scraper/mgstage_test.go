package scraper

import (
	"testing"
)

func TestMGStageScraper_FetchDoc(t *testing.T) {
	BeforeTest()
	tests := []testCase{
		{
			name: "fetchDoc expects no error",
			args: args{
				query: "abw-108",
			},
			wantErr: false,
			want:    "ABW-108",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, ok := mgsTests[tt.args.query]; !ok {
				mgsTests[tt.args.query] = &MGStageScraper{
					doc: tt.fields.doc,
				}
			}
			if err := mgsTests[tt.args.query].FetchDoc(tt.args.query); (err != nil) != tt.wantErr {
				t.Errorf("FetchDoc() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := mgsTests[tt.args.query].GetNumber(); got != tt.want {
				t.Errorf("GetNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
