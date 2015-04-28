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
	"strings"
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
	base = base[:strings.Index(base, ".")]

	for _, spec := range LaunchImageSpecifications {
		savePath := base + spec.Postfix

		log.Println("generating ", savePath)
		if cut == true {
			// if the origin is smaller than expected
			if bounds.Size().X >= spec.Width && bounds.Size().Y < spec.Height {
				// scale first, then cut
				m = resize.Resize(0, uint(spec.Height), m, resize.Bilinear)
				log.Printf("first resized to %d * %d\n", m.Bounds().Size().X, m.Bounds().Size().Y)
				goto do_cut
			}

			if bounds.Size().X < spec.Width && bounds.Size().Y >= spec.Height {
				// scale first, then cut
				m = resize.Resize(uint(spec.Width), 0, m, resize.Bilinear)
				log.Printf("or resized to %d * %d\n", m.Bounds().Size().X, m.Bounds().Size().Y)
				goto do_cut
			}

			if bounds.Size().X < spec.Width && bounds.Size().Y < spec.Height {
				// fall through, just scale
				goto do_scale
			}
		do_cut:
			// if the origin is larger than expected
			im, err := cutter.Crop(m, cutter.Config{
				Width:  spec.Width,
				Height: spec.Height,
				Anchor: image.Point{(m.Bounds().Size().X - spec.Width) / 2, (m.Bounds().Size().Y - spec.Height) / 2},
			})
			log.Printf("cropped to %d * %d\n", im.Bounds().Size().X, im.Bounds().Size().Y)

			if err != nil {
				log.Println(savePath, err)
				continue
			}

			if err := saveImage(&im, savePath, 1); err != nil {
				log.Println(savePath, err)
			}
			continue
		}
	do_scale:
		// scale it
		var im image.Image
		if bounds.Size().X > spec.Width &&
			bounds.Size().Y*spec.Width/bounds.Size().X < spec.Height {
			im = resize.Resize(uint(spec.Width), 0, m, resize.Bilinear)
		} else {
			im = resize.Resize(0, uint(spec.Height), m, resize.Bilinear)
		}
		log.Printf("finally, resized to %d * %d\n", im.Bounds().Size().X, im.Bounds().Size().Y)
		if err := saveImage(&im, savePath, 1); err != nil {
			log.Println(savePath, err)
		}
	}
	return nil
}

func GenerateAppIcon(origin string) error {
	return nil
}
