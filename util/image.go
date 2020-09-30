package util

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
	IT_png = iota
	IT_jpeg
	IT_gif
	IT_webp
	IT_tiff
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
	case IT_png:
		err = png.Encode(file, rgba)
	case IT_jpeg:
		err = jpeg.Encode(file, rgba, &jpeg.Options{Quality: 100})
	case IT_gif:
		err = gif.Encode(file, rgba, &gif.Options{})
	case IT_webp:
		err = webp.Encode(file, rgba, &webp.Options{Lossless: true})
	case IT_tiff:
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
	case IT_png:
		err = png.Encode(file, *img)
	case IT_jpeg:
		err = jpeg.Encode(file, *img, &jpeg.Options{Quality: 100})
	case IT_gif:
		err = gif.Encode(file, *img, &gif.Options{})
	case IT_webp:
		err = webp.Encode(file, *img, &webp.Options{Lossless: true})
	case IT_tiff:
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
