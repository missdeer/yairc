package main

import (
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
)

func isDir(path string) (bool, error) {
	var file *os.File
	if f, err := os.OpenFile(path, os.O_RDONLY, 0644); err != nil {
		log.Println(err)
		return false, err
	} else {
		file = f
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		log.Println(err)
		return false, err
	}
	return fi.IsDir(), nil
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
		err = jpeg.Encode(file, *img, &jpeg.Options{100})
	}

	if err != nil {
		log.Println(savePath, err)
		return err
	}
	return nil
}
