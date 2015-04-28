package main

import (
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
)

type LaunchImageSpec struct {
	Width   int
	Height  int
	Postfix string
}

var (
	LaunchImageSpecifications = []LaunchImageSpec{
		{320, 480, "~iphone.png"},
		{640, 960, "@2x~iphone.png"},
		{640, 1136, "-568h@2x~iphone.png"},
		{750, 1334, "-667h@2x~iphone.png"},
		{1242, 2208, "-736h@3x~iphone.png"},
		{768, 1024, "-Portrait~ipad.png"},
		{1024, 768, "-Landscape~ipad.png"},
		{1536, 2048, "-Portrait@2x~ipad.png"},
		{2048, 1536, "-Landscape@2x~ipad.png"},
	}
)

func GenerateLaunchImage(origin string) error {
	reader, err := os.Open(origin)
	if err != nil {
		log.Println(origin, err)
		return err
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Println(origin, err)
		return err
	}
	bounds := m.Bounds()

	base := filepath.Base(origin)

	for _, spec := range LaunchImageSpecifications {
		savePath := base + spec.Postfix
		// if the origin is larger than expected

		// if the origin is smaller than expected

		if cut == true {
			// cut first
			croppedImg, err := cutter.Crop(m, cutter.Config{
				Width:  spec.Width,
				Height: spec.Height,
				Anchor: image.Point{0, 0},
			})

			if err != nil {
				log.Println(savePath, err)
				return err
			}

			if err := saveImage(&croppedImg, savePath, 1); err != nil {
				return err
			}
		}
		if scale == true {
			// scale it
			var im image.Image
			if bounds.Size().X > spec.Width &&
				bounds.Size().Y*spec.Width/bounds.Size().X < spec.Height {
				im = resize.Resize(uint(spec.Width), 0, m, resize.Bilinear)
			} else {
				im = resize.Resize(0, uint(spec.Height), m, resize.Bilinear)
			}
			if err := saveImage(&im, savePath, 1); err != nil {
				return err
			}
		}
	}
	return nil
}

func GenerateAppIcon(origin string) error {
	return nil
}
