package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jackmordaunt/icns"
	"github.com/missdeer/yairc/util"
	flag "github.com/spf13/pflag"
)

var (
	compress            bool
	backgroundImagePath string
	foregroundImagePath string
	inputImagePath      string
	outputDirectoryPath string
	action              string
	platform                   = "ios"
	red                 uint32 = 127
	green               uint32 = 127
	blue                uint32 = 127
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
	flag.StringVarP(&action, "action", "a", "", "candidats: icons, icns, appIcon, launchImage, transparent, invert, info")
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

	// icon scale mode
	if action == "icons" && inputImagePath != "" && outputDirectoryPath != "" {
		log.Println("generate /@2x/@3x/@4x & /x18/x36/x48 icons from", inputImagePath, "to", outputDirectoryPath)
		err := iconScale(inputImagePath, outputDirectoryPath)
		if err != nil {
			log.Fatal(err)
		}
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
		pngf, err := util.OpenURI(inputImagePath)
		if err != nil {
			log.Fatalf("opening source image: %v", err)
		}
		defer pngf.Close()
		srcImg, _, err := util.ImageDecode(pngf)
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

	args := flag.Args()
	if inputImagePath != "" {
		args = append(args, inputImagePath)
	}

	if action == "transparent" && len(args) > 0 {
		log.Println("transparent color")
		for _, uri := range args {
			im, err := util.Transparent(uri, red, green, blue)
			if err != nil {
				log.Println(err)
				continue
			}
			fn := inputImagePath[:len(inputImagePath)-len(filepath.Ext(inputImagePath))] + ".transparent.png"
			if err = util.SaveImage(im, fn, util.IT_png); err != nil {
				log.Println(err)
				continue
			}
			if err = util.DoCrush(compress, fn); err != nil {
				log.Println(err)
				continue
			}
		}
		return
	}

	if action == "invert" && len(args) > 0 {
		log.Println("invert color")
		for _, uri := range args {
			im, err := util.Invert(uri)
			if err != nil {
				log.Println(err)
				continue
			}
			fn := inputImagePath[:len(inputImagePath)-len(filepath.Ext(inputImagePath))] + ".transparent.png"
			if err = util.SaveImage(im, fn, util.IT_png); err != nil {
				log.Println(err)
				continue
			}
			if err = util.DoCrush(compress, fn); err != nil {
				log.Println(err)
				continue
			}
		}
		return
	}

	if action == "info" && len(args) > 0 {
		log.Println("info color")
		for _, uri := range args {
			im, err := util.Info(uri)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(uri)
			for k, v := range im {
				fmt.Printf("%s:%d\n", k, v)
			}
			fmt.Println("===============================")
		}
		return
	}
}
