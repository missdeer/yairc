package main

import (
	"errors"
	"image"
	"image/gif"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"log"
	"os"

	"github.com/chai2010/tiff"
	"github.com/chai2010/webp"
)

const (
	it_png = iota
	it_jpeg
	it_gif
	it_webp
	it_tiff
)

// draw.DrawMask(dst, dst.Bounds(), src, image.ZP, &Circle{p, r}, image.ZP, draw.Over)


func SaveRGBA(rgba *image.RGBA, savePath string, imageType int) (err error) {
	var file *os.File
	if f, err := os.OpenFile(savePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		log.Fatal(savePath, err)
	} else {
		file = f
	}
	defer file.Close()

	switch imageType {
	case it_png:
		err = png.Encode(file, rgba)
		if err == nil && compress == true {
			err = Crush(savePath)
		}
	case it_jpeg:
		err = jpeg.Encode(file, rgba, &jpeg.Options{Quality: 100})
	case it_gif:
		err = gif.Encode(file, rgba, &gif.Options{})
	case it_webp:
		err = webp.Encode(file, rgba, &webp.Options{Lossless: true})
	case it_tiff:
		err = tiff.Encode(file, rgba, &tiff.Options{})
	default:
		err = errors.New("unsupported format")
	}

	if err != nil {
		return err
	}
	return nil
}

func SaveImage(img *image.Image, savePath string, imageType int) (err error) {
	var file *os.File
	if f, err := os.OpenFile(savePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		log.Fatal(savePath, err)
	} else {
		file = f
	}
	defer file.Close()

	switch imageType {
	case it_png:
		err = png.Encode(file, *img)
		if err == nil && compress == true {
			err = Crush(savePath)
		}
	case it_jpeg:
		err = jpeg.Encode(file, *img, &jpeg.Options{Quality: 100})
	case it_gif:
		err = gif.Encode(file, *img, &gif.Options{})
	case it_webp:
		err = webp.Encode(file, *img, &webp.Options{Lossless: true})
	case it_tiff:
		err = tiff.Encode(file, *img, &tiff.Options{})
	default:
		err = errors.New("unsupported format")
	}

	if err != nil {
		return err
	}
	return nil
}

func ImageDecode(r io.Reader) (image.Image, string, error) {
	return image.Decode(r)
}
