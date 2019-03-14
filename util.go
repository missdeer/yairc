package main

import (
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
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

//draw.DrawMask(dst, dst.Bounds(), src, image.ZP, &Circle{p, r}, image.ZP, draw.Over)

func isDir(path string) (bool, error) {
	var file *os.File
	if f, err := os.OpenFile(path, os.O_RDONLY, 0644); err != nil {
		log.Println(err)
		return false, err
	}
	file = f

	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		log.Println(err)
		return false, err
	}
	return fi.IsDir(), nil
}

func saveRGBA(rgba *image.RGBA, savePath string, imagetype int) (err error) {
	var file *os.File
	if f, err := os.OpenFile(savePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		log.Fatal(savePath, err)
	} else {
		file = f
	}
	defer file.Close()

	switch imagetype {
	case 1:
		err = png.Encode(file, rgba)
	default:
		err = jpeg.Encode(file, rgba, &jpeg.Options{Quality: 100})
	}

	if err != nil {
		log.Println(savePath, err)
		return err
	}
	return nil
}

func saveImage(img *image.Image, savePath string, imagetype int) (err error) {
	var file *os.File
	if f, err := os.OpenFile(savePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		log.Fatal(savePath, err)
	} else {
		file = f
	}
	defer file.Close()

	switch imagetype {
	case 1:
		err = png.Encode(file, *img)
	default:
		err = jpeg.Encode(file, *img, &jpeg.Options{Quality: 100})
	}

	if err != nil {
		log.Println(savePath, err)
		return err
	}
	return nil
}
