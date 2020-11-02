package util

import (
	"errors"
	"image"
	"image/draw"
	"log"
)

var (
	cutFunctionMap = map[string]cutFunc{
		"l":          cutLeft,
		"left":       cutLeft,
		"r":          cutRight,
		"right":      cutRight,
		"b":          cutBottom,
		"bottom":     cutBottom,
		"v":          cutVertical,
		"vertical":   cutVertical,
		"h":          cutHorizontal,
		"horizontal": cutHorizontal,
		"a":          cutAll,
		"all":        cutAll,
	}
)

type cutFunc func(image.Image, uint) (image.Image, error)

func cutLeft(im image.Image, step uint) (image.Image, error) {
	rc := im.Bounds()
	bm := image.NewRGBA(image.Rect(0, 0, rc.Dx()-int(step), rc.Dy()))
	draw.Draw(bm, bm.Bounds(), im, image.Point{int(step), 0}, draw.Over)
	return bm, nil
}

func cutRight(im image.Image, step uint) (image.Image, error) {
	rc := im.Bounds()
	bm := image.NewRGBA(image.Rect(0, 0, rc.Dx()-int(step), rc.Dy()))
	draw.Draw(bm, bm.Bounds(), im, image.Point{0, 0}, draw.Over)
	return bm, nil
}

func cutTop(im image.Image, step uint) (image.Image, error) {
	rc := im.Bounds()
	bm := image.NewRGBA(image.Rect(0, 0, rc.Dx(), rc.Dy()-int(step)))
	draw.Draw(bm, bm.Bounds(), im, image.Point{0, int(step)}, draw.Over)
	return bm, nil
}

func cutBottom(im image.Image, step uint) (image.Image, error) {
	rc := im.Bounds()
	bm := image.NewRGBA(image.Rect(0, 0, rc.Dx(), rc.Dy()-int(step)))
	draw.Draw(bm, bm.Bounds(), im, image.Point{0, 0}, draw.Over)
	return bm, nil
}

func cutHorizontal(im image.Image, step uint) (image.Image, error) {
	i, err := cutRight(im, step)
	if err != nil {
		return nil, err
	}
	return cutLeft(i, step)
}

func cutVertical(im image.Image, step uint) (image.Image, error) {
	i, err := cutTop(im, step)
	if err != nil {
		return nil, err
	}
	return cutBottom(i, step)
}

func cutAll(im image.Image, step uint) (image.Image, error) {
	i, err := cutHorizontal(im, step)
	if err != nil {
		return nil, err
	}
	return cutVertical(i, step)
}

func CutEdge(uri string, position string, step uint) (image.Image, error) {
	f, ok := cutFunctionMap[position]
	if !ok {
		return nil, errors.New("invalid cut edge position")
	}

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
	if rc.Dy() <= int(step) || rc.Dx() <= int(step) {
		return nil, errors.New("invalid cut edge step")
	}

	return f(im, step)
}
