package main

import (
	"log"
	"os"

	"github.com/nfnt/resize"
)

type SplashScreenSpec struct {
	Width  int
	Height int
	Name   string
}

type LauncherIconSpec struct {
	Length int
	Name   string
}

var (
	SplashScreenSpecifications = []SplashScreenSpec{
		{768, 1280, "drawable-xhdpi"},
		{800, 1280, "drawable-213dpi"},
		{1080, 1920, "drawable-xxhdpi"},
		{1440, 2560, "drawable-560dpi"},
	}
	LauncherIconSpecifications = []LauncherIconSpec{
		{36, "drawable-ldpi"},
		{48, "drawable-mdpi"},
		{64, "drawable-tvdpi"},
		{72, "drawable-hdpi"},
		{96, "drawable-xhdpi"},
		{144, "drawable-xxhdpi"},
		{192, "darwable-xxxhdpi"},
	}
)

func GenerateSplashScreen(origin string) error {
	return nil
}

func GenerateLauncherIcon(origin string) error {
	reader, err := os.Open(origin)
	if err != nil {
		log.Println(origin, err)
		return err
	}
	defer reader.Close()
	m, _, err := ImageDecode(reader)
	if err != nil {
		log.Println(origin, err)
		return err
	}

	for _, spec := range LauncherIconSpecifications {
		im := resize.Resize(uint(spec.Length), uint(spec.Length), m, resize.Bilinear)
		if err := saveImage(&im, spec.Name, it_png); err != nil {
			log.Println(spec.Name, err)
		}
	}
	return nil
}
