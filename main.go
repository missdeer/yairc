package main

import (
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	MaxWidth  int = 620
	MaxHeight int = 960
)

func scale(filepath string, imagetype int) error {
	reader, err := os.Open(filepath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Println(err)
		return err
	}
	bounds := m.Bounds()

	if bounds.Size().X > MaxWidth || bounds.Size().Y > MaxHeight {
		// scale it
		var im image.Image
		if bounds.Size().X > MaxWidth && bounds.Size().Y*MaxWidth/bounds.Size().X < MaxHeight {
			im = resize.Resize(uint(MaxWidth), 0, m, resize.Bilinear)
		} else {
			im = resize.Resize(0, uint(MaxHeight), m, resize.Bilinear)
		}
		var savePath string
		switch imagetype {
		case 1:
			savePath = filepath + "-m.png"
		default:
			savePath = filepath + "-m.jpg"
		}

		fmt.Printf("%s width=%d, height=%d, saved to %s \n", filepath, bounds.Size().X, bounds.Size().Y, savePath)

		var file *os.File
		if f, err := os.OpenFile(savePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
			log.Fatal(err)
		} else {
			file = f
		}
		defer file.Close()

		switch imagetype {
		case 1:
			err = png.Encode(file, im)
		default:
			err = jpeg.Encode(file, im, &jpeg.Options{100})
		}

		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		if strings.LastIndex(path, "-m.jpg") < 0 && strings.LastIndex(path, "-m.png") < 0 {
			savePath := path + "-m.jpg"
			if _, err := os.Stat(savePath); err != nil {
				// not exists
				scale(path, 2)
			}

			savePath = path + "-m.png"
			if _, err := os.Stat(savePath); err != nil {
				// not exists
				scale(path, 1)
			}
		}
	} else {
		fmt.Printf("Visited: %s\n", path)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("Incorrect arguments!"))
	}
	for _, root := range os.Args {
		var file *os.File
		if f, err := os.OpenFile(root, os.O_RDONLY, 0644); err != nil {
			log.Println(err)
			continue
		} else {
			file = f
		}
		defer file.Close()
		fi, err := file.Stat()
		if err != nil {
			log.Println(err)
			continue
		}
		if fi.IsDir() {
			err := filepath.Walk(root, visit)
			if err != nil {
				log.Println(err)
				continue
			}
		} else {
			scale(root, 1)
			scale(root, 2)
		}
	}
}
