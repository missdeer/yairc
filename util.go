package main

import (
	"image"
	"image/color"
	"image/gif"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"log"
	"os"

	"github.com/chai2010/webp"
)

const (
	it_png = iota
	it_jpeg
	it_gif
	it_webp
)

type Circle struct {
	p image.Point
	r int
}

func (c *Circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *Circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *Circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

// draw.DrawMask(dst, dst.Bounds(), src, image.ZP, &Circle{p, r}, image.ZP, draw.Over)

func isDir(path string) (bool, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		log.Println(err)
		return false, err
	}
	file := f

	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		log.Println(err)
		return false, err
	}
	return fi.IsDir(), nil
}

func saveRGBA(rgba *image.RGBA, savePath string, imageType int) (err error) {
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
			err = crush(savePath)
		}
	case it_jpeg:
		err = jpeg.Encode(file, rgba, &jpeg.Options{Quality: 100})
	case it_gif:
		err = gif.Encode(file, rgba, &gif.Options{})
	case it_webp:
		err = webp.Encode(file, rgba, &webp.Options{Lossless: true})
	default:
	}

	if err != nil {
		log.Println(savePath, err)
		return err
	}
	return nil
}

func saveImage(img *image.Image, savePath string, imageType int) (err error) {
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
			err = crush(savePath)
		}
	case it_jpeg:
		err = jpeg.Encode(file, *img, &jpeg.Options{Quality: 100})
	case it_gif:
		err = gif.Encode(file, *img, &gif.Options{})
	case it_webp:
		err = webp.Encode(file, *img, &webp.Options{Lossless: true})
	default:
	}

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func ImageDecode(r io.Reader) (image.Image, string, error) {
	m, err := webp.Decode(r)
	if err != nil {
		return image.Decode(r)
	}
	return m, "webp", err
}
