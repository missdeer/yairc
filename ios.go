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

type handler func(image.Image, string, *LaunchImageSpec) error

type LaunchImageSpec struct {
	Width   int
	Height  int
	Postfix string
	Handler handler
}

var (
	LaunchImageSpecifications = []LaunchImageSpec{
		{320, 480, "~iphone.png", ScaleHandler},
		{640, 960, "@2x~iphone.png", ScaleHandler},
		{640, 1136, "-568h@2x~iphone.png", ScaleCutHandler},
		{750, 1334, "-667h@2x~iphone.png", ScaleCutHandler},
		{1242, 2208, "-736h@3x~iphone.png", ScaleCutHandler},
		{768, 1024, "-Portrait~ipad.png", ScaleHandler},
		{1024, 768, "-Landscape~ipad.png", ScaleCutHandler},
		{1536, 2048, "-Portrait@2x~ipad.png", SkipHandler},
		{2048, 1536, "-Landscape@2x~ipad.png", ScaleCutHandler},
	}
)

func ScaleHandler(m image.Image, savePath string, spec *LaunchImageSpec) error {
	bounds := m.Bounds()
	var im image.Image

	if bounds.Size().Y*spec.Width/bounds.Size().X < spec.Height {
		im = resize.Resize(uint(spec.Width), 0, m, resize.Bilinear)
		log.Println("resize by width")
	} else {
		im = resize.Resize(0, uint(spec.Height), m, resize.Bilinear)
		log.Println("resize by height")
	}

	if err := saveImage(&im, savePath, 1); err != nil {
		log.Println(savePath, err)
	}
	return nil
}

func ScaleCutHandler(m image.Image, savePath string, spec *LaunchImageSpec) error {
	bounds := m.Bounds()
	var im image.Image

	if bounds.Size().Y*spec.Width/bounds.Size().X < spec.Height {
		im = resize.Resize(uint(spec.Width), 0, m, resize.Bilinear)
		log.Println("resize by width")
	} else {
		im = resize.Resize(0, uint(spec.Height), m, resize.Bilinear)
		log.Println("resize by height")
	}

	// if the origin is larger than expected
	var err error
	im, err = cutter.Crop(im, cutter.Config{
		Width:  spec.Width,
		Height: spec.Height,
		Anchor: image.Point{(m.Bounds().Size().X - spec.Width) / 2, (m.Bounds().Size().Y - spec.Height) / 2},
	})
	log.Printf("cropped to %d * %d\n", im.Bounds().Size().X, im.Bounds().Size().Y)

	if err != nil {
		log.Println(savePath, err)
		return err
	}
	if err := saveImage(&im, savePath, 1); err != nil {
		log.Println(savePath, err)
		return err
	}
	return nil
}

func SkipHandler(m image.Image, savePath string, spec *LaunchImageSpec) error {
	if err := saveImage(&m, savePath, 1); err != nil {
		log.Println(savePath, err)
	}
	return nil
}

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

	base := filepath.Base(origin)
	base = base[:strings.Index(base, ".")]

	for _, spec := range LaunchImageSpecifications {
		savePath := base + spec.Postfix

		log.Println("generating ", savePath)
		spec.Handler(m, savePath, &spec)
	}
	return nil
}

func GenerateAppIcon(origin string) error {
	return nil
}
