package grabber

import (
	"better-av-tool/log"
	"github.com/oliamb/cutter"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

const posterWidth = 378

func Crop(src, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		log.Error("Cannot open file", err)
		return err
	}
	defer f.Close()
	img, err := jpeg.Decode(f)
	if err != nil {
		log.Error("Cannot decode image:", err)
		return err
	}
	srcW := img.Bounds().Dx()
	srcH := img.Bounds().Dy()
	img, err = cutter.Crop(img, cutter.Config{
		Height:  srcH,                                     // height in pixel or Y ratio(see Ratio Option below)
		Width:   posterWidth,                              // width in pixel or X ratio
		Mode:    cutter.TopLeft,                           // Accepted Mode: TopLeft, Centered
		Anchor:  image.Point{X: srcW - posterWidth, Y: 0}, // Position of the top left point
		Options: 0,                                        // Accepted Option: Ratio
	})

	if err != nil {
		log.Error("Cannot Crop image:", err)
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		log.Error("Cannot create image:", err)
		return err
	}
	defer out.Close()
	err = png.Encode(out, img)

	if err != nil {
		log.Error("Cannot Encode image:", err)
		return err
	}

	return nil
}
