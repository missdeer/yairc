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

type AppIconSpec struct {
	Length int
	Name   string
}

var (
	LaunchImageSpecifications = []LaunchImageSpec{
		{320, 480, "~iphone.png", ScaleCutHandler},
		{640, 960, "@2x~iphone.png", ScaleCutHandler},
		{640, 1136, "-568h@2x~iphone.png", ScaleCutHandler},
		{750, 1334, "-667h@2x~iphone.png", ScaleCutHandler},
		{1242, 2208, "-736h@3x~iphone.png", ScaleCutHandler},
		{768, 1024, "-Portrait~ipad.png", ScaleHandler},
		{1024, 768, "-Landscape~ipad.png", ScaleCutHandler},
		{1536, 2048, "-Portrait@2x~ipad.png", SkipHandler},
		{2048, 1536, "-Landscape@2x~ipad.png", ScaleCutHandler},
	}
	AppIconSpecifications = []AppIconSpec{
		{20, "Icon-Small-20.png"},
		{29, "Icon-29.png"},
		{29, "Icon-Small.png"},
		{30, "Icon-Small-30.png"},
		{40, "Icon-40.png"},
		{40, "Icon-Small-40.png"},
		{40, "Icon-Small-20@2x.png"},
		{50, "Icon-Small-50.png"},
		{57, "Icon.png"},
		{58, "Icon-29@2x.png"},
		{58, "Icon-Small@2x.png"},
		{60, "Icon-Small-30@2x.png"},
		{72, "Icon-72.png"},
		{76, "Icon-76.png"},
		{80, "Icon-40@2x.png"},
		{80, "Icon-Small-40@2x.png"},
		{87, "Icon-29@3x.png"},
		{100, "Icon-Small-50@2x.png"},
		{114, "Icon@2x.png"},
		{120, "Icon-40@3x.png"},
		{120, "Icon-60@2x.png"},
		{120, "Icon-120.png"},
		{144, "Icon-72@2x.png"},
		{152, "Icon-76@2x.png"},
		{180, "Icon-60@3x.png"},
		{512, "iTunesArtwork.png"},
		{1024, "iTunesArtwork@2x.png"},
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

	for _, spec := range AppIconSpecifications {
		im := resize.Resize(uint(spec.Length), uint(spec.Length), m, resize.Bilinear)
		if err := saveImage(&im, spec.Name, 1); err != nil {
			log.Println(spec.Name, err)
		}
	}
	return nil
}
