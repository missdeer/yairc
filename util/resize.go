package util

import (
	"image"
	"log"

	"github.com/nfnt/resize"
)

func Resize(uri string, w, h uint) (image.Image, error) {
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

	return resize.Resize(w, h, im, resize.Bilinear), nil
}
