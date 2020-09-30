package util

import (
	"errors"
	"image/color"
	"image/draw"
	"log"
)

func Info(uri string) (map[color.Color]int, error) {
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
	cm := make(map[color.Color]int)
	for x := 0; x < rc.Dx(); x++ {
		for y := 0; y < rc.Dy(); y++ {
			c := img.At(x, y)
			cm[c]++
		}
	}
	return cm, nil
}
