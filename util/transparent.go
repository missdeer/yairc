package util

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"log"
)

func Transparent(uri string, red, green, blue uint32, transparentWhiteDirect bool) (image.Image, error) {
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
			r, g, b, _ := c.RGBA()

			if !transparentWhiteDirect && r > red && g > green && b > blue {
				img.Set(x, y, color.Transparent)
			}
			if transparentWhiteDirect && r < red && g < green && b < blue {
				img.Set(x, y, color.Transparent)
			}
		}
	}
	return im, nil
}
