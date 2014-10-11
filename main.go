// This example demonstrates decoding a JPEG image and examining its pixels.
package main

import (
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
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

func scale(filepath string) error {
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
		savePath := filepath + "-m.png"

		fmt.Printf("%s width=%d, height=%d, saved to %s \n", filepath, bounds.Size().X, bounds.Size().Y, savePath)

		var file *os.File
		if f, err := os.OpenFile(savePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
			log.Fatal(err)
		} else {
			file = f
		}
		defer file.Close()

		err = png.Encode(file, im)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		if strings.LastIndex(path, "-m.png") < 0 {
			scale(path)
		}
	} else {
		fmt.Printf("Visited: %s\n", path)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal(errors.New("Incorrect arguments!"))
	}
	root := os.Args[1]

	var file *os.File
	if f, err := os.OpenFile(root, os.O_RDONLY, 0644); err != nil {
		log.Fatal(err)
	} else {
		file = f
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if fi.IsDir() {
		err := filepath.Walk(root, visit)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		scale(root)
	}
}
