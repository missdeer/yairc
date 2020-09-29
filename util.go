package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"errors"
	"image"
	"image/color"
	"image/gif"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/andybalholm/brotli"
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
	case it_tiff:
		err = tiff.Encode(file, rgba, &tiff.Options{})
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
	case it_tiff:
		err = tiff.Encode(file, *img, &tiff.Options{})
	default:
	}

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func ImageDecode(r io.Reader) (image.Image, string, error) {
	return image.Decode(r)
}

func uncompressReader(r *http.Response) (io.ReadCloser, bool, error) {
	header := strings.ToLower(r.Header.Get("Content-Encoding"))
	switch header {
	case "":
		return r.Body, false, nil
	case "br":
		rc := brotli.NewReader(r.Body)
		if rc == nil {
			log.Println("creating brotli reader failed")
			return nil, false, errors.New("creating brotli reader failed")
		}
		return ioutil.NopCloser(rc), true, nil
	case "gzip":
		rc, err := gzip.NewReader(r.Body)
		if err != nil {
			log.Println("creating gzip reader failed:", err)
			return nil, false, err
		}
		return rc, true, nil
	case "deflate":
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("reading inflate failed:", err)
			return nil, false, err
		}
		rc := flate.NewReader(bytes.NewReader(content[2:]))
		if rc == nil {
			log.Println("creating deflate reader failed")
			return nil, false, errors.New("creating deflate reader failed")
		}
		return rc, true, nil
	}
	return nil, false, errors.New("unexpected encoding type")
}

func OpenURI(uri string) (rc io.ReadCloser, err error) {
	if strings.HasPrefix(uri, "https://") || strings.HasPrefix(uri, "http://") {
		req, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			return os.Open(uri)
		}

		req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")

		httpClient := http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			return os.Open(uri)
		}

		if resp.StatusCode != 200 {
			return os.Open(uri)
		}

		rc, _, err = uncompressReader(resp)
		if err == nil {
			return rc, err
		}
	}

	return os.Open(uri)
}
