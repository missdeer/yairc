package util

import (
	"image"
	"image/color"
	"image/draw"
	"log"
)

// enforce image.RGBA to always add the alpha channel when encoding PNGs.
type notOpaqueRGBA struct {
	*image.RGBA
}

func (i *notOpaqueRGBA) Opaque() bool {
	return false
}

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

	rc := im.Bounds()
	// https://groob.io/posts/image-draw-intro/
	img := &notOpaqueRGBA{image.NewRGBA(im.Bounds())}
	draw.Draw(img, rc, im, image.ZP, draw.Src)

	for x := 0; x < rc.Dx(); x++ {
		for y := 0; y < rc.Dy(); y++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			r &= 0xff
			g &= 0xff
			b &= 0xff
			if !transparentWhiteDirect && r > red && g > green && b > blue {
				img.Set(x, y, color.Transparent)
			}
			if transparentWhiteDirect && r < red && g < green && b < blue {
				img.Set(x, y, color.Transparent)
			}
		}
	}
	return img, nil
}
