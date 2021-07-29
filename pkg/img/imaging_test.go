package img

import (
	"testing"
)

func TestImaging_CropAndSave(t *testing.T) {
	type args struct {
		src string
		dst string
		w   int
		h   int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test crop",
			args: args{
				src: "test.jpg",
				dst: "test1.jpg",
				w:   378,
				h:   0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Imaging{}
			if err := i.CropAndSave(tt.args.src, tt.args.dst, tt.args.w, tt.args.h); (err != nil) != tt.wantErr {
				t.Errorf("CropAndSave() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
