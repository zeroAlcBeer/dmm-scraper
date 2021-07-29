package img

import (
	"image"

	"github.com/disintegration/imaging"
)

// Imaging ...
type Imaging struct {
}

// Crop ...
func (i *Imaging) Crop(src image.Image, w, h int) (image.Image, error) {
	return imaging.CropAnchor(src, w, h, imaging.TopRight), nil
}

// Open ...
func (i *Imaging) Open(filename string) (image.Image, error) {
	return imaging.Open(filename)
}

// Save ...
func (i *Imaging) Save(img image.Image, filename string) (err error) {
	return imaging.Save(img, filename)
}

// CropAndSave ...
func (i *Imaging) CropAndSave(src, dst string, w, h int) error {
	img, err := i.Open(src)
	if err != nil {
		return err
	}
	if w == 0 {
		w = img.Bounds().Dx()
	}
	if h == 0 {
		h = img.Bounds().Dy()
	}
	croped, err := i.Crop(img, w, h)
	if err != nil {
		return err
	}
	return i.Save(croped, dst)
}
