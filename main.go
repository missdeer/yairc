package main

import (
	"fmt"
	"image/color"
	"image/draw"
	"log"
	"os"
	"path/filepath"

	"github.com/jackmordaunt/icns"
	flag "github.com/spf13/pflag"
)

const (
	MaxWidth  int = 620
	MaxHeight int = 960
)

var (
	compress            bool
	imageDirectory      string
	backgroundImagePath string
	foregroundImagePath string
	inputImagePath      string
	outputDirectoryPath string
	action              string
	platform            = "ios"

	red   uint32 = 127
	green uint32 = 127
	blue  uint32 = 127
	// Gitcommit contains the commit where we built from.
	GitCommit string
)

func main() {
	showHelpMessage := false
	showVersion := false
	flag.Uint32VarP(&red, "red", "", red, "set red threshold")
	flag.Uint32VarP(&green, "green", "", green, "set green threshold")
	flag.Uint32VarP(&blue, "blue", "", blue, "set blue threshold")
	flag.BoolVarP(&compress, "compress", "", true, "compress output PNG files")
	flag.StringVarP(&platform, "platform", "p", "common", "candidates: ios, android, common")
	flag.StringVarP(&action, "action", "a", "", "candidats: icons, icns, appIcon, launchImage, transparent")
	flag.StringVarP(&backgroundImagePath, "background", "b", "", "path of background image for launch image")
	flag.StringVarP(&foregroundImagePath, "foreground", "f", "", "path of foreground image for launch image")
	flag.StringVarP(&inputImagePath, "input", "i", "", "input image path")
	flag.StringVarP(&outputDirectoryPath, "output", "o", ".", "output directory path")
	flag.BoolVarP(&showHelpMessage, "help", "h", false, "show this help message")
	flag.BoolVarP(&showVersion, "version", "v", false, "show version number")
	flag.Parse()

	if showHelpMessage {
		flag.PrintDefaults()
		return
	}

	if showVersion {
		fmt.Println("yairc version:", GitCommit)
		return
	}

	if action == "transparent" && inputImagePath != "" {
		log.Println("transparent color")
		r, err := OpenURI(inputImagePath)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()
		im, format, err := ImageDecode(r)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("found format:", format)
		img, ok := im.(draw.Image)
		if !ok {
			log.Fatal("not a drawable image type")
		}
		rc := img.Bounds()
		cm := make(map[color.Color]int)
		for x := 0; x < rc.Dx(); x++ {
			for y := 0; y < rc.Dy(); y++ {
				c := img.At(x, y)
				r, g, b, _ := c.RGBA()
				if r > red && g > green && b > blue {
					img.Set(x, y, color.Transparent)
				}
				cm[c]++
			}
		}
		for clr, count := range cm {
			fmt.Println(clr, count)
		}
		fn := inputImagePath[:len(inputImagePath)-len(filepath.Ext(inputImagePath))] + ".transparent.png"
		err = SaveImage(&im, fn, it_png)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// icon scale mode
	if action == "icons" && inputImagePath != "" && outputDirectoryPath != "" {
		log.Println("generate /@2x/@3x/@4x & /x18/x36/x48 icons from", inputImagePath, "to", outputDirectoryPath)
		iconScale(inputImagePath, outputDirectoryPath)
		return
	}

	// ios app icon mode
	if action == "appIcon" && platform == "ios" {
		fmt.Println("output ios app icons")
		err := GenerateAppIcon(inputImagePath)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// ios launch image mode
	if action == "launchImage" && platform == "ios" {
		fmt.Println("output ios launch images")
		err := GenerateLaunchImage()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if action == "appIcon" && platform == "android" {
		fmt.Println("output android launcher icons")
		err := GenerateLauncherIcon(inputImagePath)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if action == "launchImage" && platform == "ios" {
		fmt.Println("output android splash screen images")
		err := GenerateSplashScreen(inputImagePath)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// convert to .icns file
	if action == "icns" && inputImagePath != "" {
		pngf, err := OpenURI(inputImagePath)
		if err != nil {
			log.Fatalf("opening source image: %v", err)
		}
		defer pngf.Close()
		srcImg, _, err := ImageDecode(pngf)
		if err != nil {
			log.Fatalf("decoding source image: %v", err)
		}
		ext := filepath.Ext(inputImagePath)
		outfile := inputImagePath[0:len(inputImagePath)-len(ext)] + ".icns"
		dest, err := os.Create(outfile)
		if err != nil {
			log.Fatalf("opening destination file: %v", err)
		}
		defer dest.Close()
		if err := icns.Encode(dest, srcImg); err != nil {
			log.Fatalf("encoding icns: %v", err)
		}
		return
	}

}
