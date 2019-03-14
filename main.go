package main

import (
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jackmordaunt/icns"
	flag "github.com/ogier/pflag"
)

const (
	MaxWidth  int = 620
	MaxHeight int = 960
)

var (
	watcher              *fsnotify.Watcher
	cut                  bool
	scale                bool
	iconScaleMode        bool
	iosScale             bool
	iosScaleTemplateSize string
	iosAppIcon           bool
	iosLaunchImage       bool
	androidLauncherIcon  bool
	androidSplash        bool
	watch                bool
	icnsMode             bool
	imagedirectory       string
	backgroundImagePath  string
	foregroundImagePath  string
	inputImagePath       string
	outputDirectoryPath  string
)

func main() {
	flag.BoolVarP(&cut, "cut", "c", true, "cut mode")
	flag.BoolVarP(&scale, "scale", "s", false, "scale mode")
	flag.BoolVarP(&iconScaleMode, "iconScale", "", false, "generate _/@2x/@3x/@4x & _/x18/x36/x48 icons")
	flag.BoolVarP(&iosScale, "iOSScale", "", false, "generate @1x/@2x/@3x images for iOS")
	flag.StringVarP(&iosScaleTemplateSize, "as", "", "1x", "iOS scale mode template size, can be 1x/2x/3x")
	flag.BoolVarP(&iosAppIcon, "appIcon", "a", false, "generate ios app icons")
	flag.BoolVarP(&iosLaunchImage, "launchImage", "l", false, "generate ios launch images")
	flag.BoolVarP(&androidLauncherIcon, "launcherIcon", "u", false, "generate android launcher icons")
	flag.BoolVarP(&androidSplash, "splashScreen", "r", false, "generate android splash screen images")
	flag.BoolVarP(&watch, "watch", "w", false, "watch directories change")
	flag.StringVarP(&backgroundImagePath, "background", "b", "", "path of background image for launch image")
	flag.StringVarP(&foregroundImagePath, "foreground", "f", "", "path of foreground image for launch image")
	flag.StringVarP(&inputImagePath, "input", "i", "", "input image path")
	flag.StringVarP(&outputDirectoryPath, "output", "o", "", "output directory path")
	flag.BoolVarP(&icnsMode, "icns", "", false, "convert input image file to .icns file")
	flag.Parse()
	if len(os.Args) < 2 {
		log.Fatal("Incorrect arguments! Use --help to get the usage.")
	}
	cut = !scale
	args := flag.Args()

	// icon scale mode
	if iconScaleMode == true && inputImagePath != "" && outputDirectoryPath != "" {
		fmt.Println("generate /@2x/@3x/@4x & /x18/x36/x48 icons")
		iconScale(inputImagePath, outputDirectoryPath)
		return
	}

	// ios scale Mode
	if iosScale == true {
		fmt.Println("generate @1x/@2x/@3x images for iOS")
		for _, root := range args {
			iOSScale(root, iosScaleTemplateSize)
		}
		return
	}

	// ios app icon mode
	if iosAppIcon == true {
		fmt.Println("output ios app icons")
		for _, root := range args {
			GenerateAppIcon(root)
		}
		return
	}

	// ios launch image mode
	if iosLaunchImage == true {
		GenerateLaunchImage()
		return
	}

	// convert to .icns file
	if icnsMode && inputImagePath != "" {
		pngf, err := os.Open(inputImagePath)
		if err != nil {
			log.Fatalf("opening source image: %v", err)
		}
		defer pngf.Close()
		srcImg, _, err := image.Decode(pngf)
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

	if androidLauncherIcon == true {
		fmt.Println("output android launcher icons")
		for _, root := range args {
			GenerateLauncherIcon(root)
		}
		return
	}

	if androidSplash == true {
		fmt.Println("output android splash screen images")
		for _, root := range args {
			GenerateSplashScreen(root)
		}
	}

	// taobao mode
	if watch == true {
		var err error
		watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(errors.New("creating filesystem watcher failed"))
		}

		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if b, e := isDir(event.Name); e == nil && b == false {
						if (event.Op&fsnotify.Remove != 0) || (event.Op&fsnotify.Write != 0) {
							// delete associated files
							if strings.LastIndex(event.Name, "-m.jpg") < 0 &&
								strings.LastIndex(event.Name, "-m.png") < 0 {
								os.Remove(event.Name + "-m.jpg")
								os.Remove(event.Name + "-m.png")
							}
						}
						if event.Op&fsnotify.Write != 0 {
							doScaleImage(event.Name)
						}
					}
				case err := <-watcher.Errors:
					log.Println("error:", err)
				}
			}
		}()

		for _, root := range args {
			if b, e := isDir(root); e == nil && b == true {
				err := filepath.Walk(root, watchDirectory)
				if err != nil {
					log.Println(err)
					continue
				}
			} else {
				doScaleImage(root)
			}
		}
		fmt.Println("watching directories...")

		timer := time.NewTicker(1 * time.Hour)
		for {
			select {
			case <-timer.C:
				fmt.Println("now: ", time.Now().UTC())
			}
		}
		timer.Stop()
	}

	if cut == true {
		for _, root := range args {
			if b, e := isDir(root); e == nil && b == true {
				err := filepath.Walk(root, traverseCut)
				if err != nil {
					log.Println(err)
					continue
				}
			} else {
				doCutImage(root)
			}
		}
		return
	}

	for _, root := range args {
		if b, e := isDir(root); e == nil && b == true {
			err := filepath.Walk(root, traverseScale)
			if err != nil {
				log.Println(err)
				continue
			}
		} else {
			doScaleImage(root)
		}
	}

}
