package scraper

import "testing"

func TestGetLabelNumber(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				s: "1sdde552",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetLabelNumber(tt.args.s)
			t.Logf(got, got1)
		})
	}
}
