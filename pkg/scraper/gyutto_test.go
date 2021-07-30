package scraper

import (
	"testing"
)

func TestGyuttoScraper_FetchDoc(t *testing.T) {
	//Setup()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "fetchDoc expects no error",
			args: args{
				query: "sex friend 54",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, ok := gyuttoTests[tt.args.query]; !ok {
				gyuttoTests[tt.args.query] = &GyuttoScraper{}
			}
			if err := gyuttoTests[tt.args.query].FetchDoc(tt.args.query); (err != nil) != tt.wantErr {
				t.Errorf("FetchDoc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	TestGyuttoScraper_GetTitle(t)
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
				query: "sex friend 54",
			},
			want: "Sex Friend 54「OGF Vol.5 刑◯姫編",
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
