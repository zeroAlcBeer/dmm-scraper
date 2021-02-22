package img

import "image"

// Operation ...
type Operation interface {
	Crop(img image.Image, w, h int) (image.Image, error)
	Open(filename string) (image.Image, error)
	Save(img image.Image, filename string) error
	CropAndSave(src, dst string, w, h int) error
}

// NewOperation ...
func NewOperation() Operation {
	return &Imaging{}
}
