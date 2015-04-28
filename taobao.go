package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strings"
)

func traverseCut(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		doCutImage(path)
	} else {
		fmt.Printf("Visited: %s\n", path)
	}
	return nil
}

func doCutImage(path string) {
	if strings.LastIndex(path, ")-p.jpg") < 0 &&
		strings.LastIndex(path, ")-p.png") < 0 &&
		strings.LastIndex(path, "-m.jpg") < 0 &&
		strings.LastIndex(path, "-m.png") < 0 {
		savePath := path + "(1)-p.jpg"
		if _, err := os.Stat(savePath); err != nil {
			// not exists
			cutImage(path, 2)
		}

		savePath = path + "(1)-p.png"
		if _, err := os.Stat(savePath); err != nil {
			// not exists
			cutImage(path, 1)
		}
	}
}

func cutImage(filepath string, imagetype int) error {
	reader, err := os.Open(filepath)
	if err != nil {
		log.Println(filepath, err)
		return err
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Println(filepath, err)
		return err
	}
	bounds := m.Bounds()

	expectHeight := bounds.Size().X * MaxHeight / MaxWidth
	if bounds.Size().Y > expectHeight {
		blockCount := bounds.Size().Y/expectHeight + 1

		for i := 0; i < blockCount && i*expectHeight < bounds.Size().Y; i++ {
			// cut first
			croppedImg, err := cutter.Crop(m, cutter.Config{
				Width:  bounds.Size().X,
				Height: expectHeight,
				Anchor: image.Point{0, i * expectHeight},
			})

			if err != nil {
				log.Println(filepath, err)
				return err
			}

			// resize then
			im := resize.Resize(uint(MaxWidth), 0, croppedImg, resize.Bilinear)

			// save to file finally
			var savePath string
			switch imagetype {
			case 1:
				savePath = fmt.Sprintf("%s(%d)-p.png", filepath, i+1)
			default:
				savePath = fmt.Sprintf("%s(%d)-p.jpg", filepath, i+1)
			}
			fmt.Printf("%s width=%d, height=%d, cropped to %s\n",
				filepath, bounds.Size().X, bounds.Size().Y, savePath)

			if err := saveImage(&im, savePath, imagetype); err != nil {
				return err
			}
		}
	}
	return nil
}

func scaleImage(filepath string, imagetype int) error {
	reader, err := os.Open(filepath)
	if err != nil {
		log.Println(filepath, err)
		return err
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Println(filepath, err)
		return err
	}
	bounds := m.Bounds()

	if bounds.Size().X > MaxWidth || bounds.Size().Y > MaxHeight {
		// scale it
		var im image.Image
		if bounds.Size().X > MaxWidth &&
			bounds.Size().Y*MaxWidth/bounds.Size().X < MaxHeight {
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

		fmt.Printf("%s width=%d, height=%d, scaled to %s\n",
			filepath, bounds.Size().X, bounds.Size().Y, savePath)

		return saveImage(&im, savePath, imagetype)
	}

	return nil
}

func traverseScale(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		doScaleImage(path)
	} else {
		fmt.Printf("Visited: %s\n", path)
	}
	return nil
}

func doScaleImage(path string) {
	if strings.LastIndex(path, "-m.jpg") < 0 &&
		strings.LastIndex(path, "-m.png") < 0 {
		savePath := path + "-m.jpg"
		if _, err := os.Stat(savePath); err != nil {
			// not exists
			scaleImage(path, 2)
		}

		savePath = path + "-m.png"
		if _, err := os.Stat(savePath); err != nil {
			// not exists
			scaleImage(path, 1)
		}
	}
}

func watchDirectory(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		if err := watcher.Add(path); err != nil {
			log.Println(err)
		}
	} else {
		doScaleImage(path)
	}
	return nil
}
