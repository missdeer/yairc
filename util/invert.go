package util

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"log"
)

func Invert(uri string) (image.Image, error) {
	r, err := OpenURI(uri)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	im, format, err := ImageDecode(r)
	if err != nil {
		return nil, err
	}
	log.Println("found format:", format)

	img, ok := im.(draw.Image)
	if !ok {
		return nil, errors.New("not a drawable image type")
	}
	rc := img.Bounds()
	for x := 0; x < rc.Dx(); x++ {
		for y := 0; y < rc.Dy(); y++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()
			img.Set(x, y, color.RGBA{uint8(^r), uint8(^g), uint8(^b), uint8(a)})
		}
	}
	return im, nil
}
