package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-fsnotify/fsnotify"
	flag "github.com/ogier/pflag"
)

const (
	MaxWidth  int = 620
	MaxHeight int = 960
)

var (
	watcher             *fsnotify.Watcher
	cut                 bool
	scale               bool
	iosAppIcon          bool
	iosLaunchImage      bool
	androidLauncherIcon bool
	androidSplash       bool
	watch               bool
	imagedirectory      string
)

func main() {
	flag.BoolVarP(&cut, "cut", "c", true, "cut mode")
	flag.BoolVarP(&scale, "scale", "s", false, "scale mode")
	flag.BoolVarP(&iosAppIcon, "appicon", "a", false, "generate ios app icons")
	flag.BoolVarP(&iosLaunchImage, "launchimage", "l", false, "generate ios launch images")
	flag.BoolVarP(&androidLauncherIcon, "launchericon", "u", false, "generate android launcher icons")
	flag.BoolVarP(&androidSplash, "splashscreen", "r", false, "generate android splash screen images")
	flag.BoolVarP(&watch, "watch", "w", false, "watch directories change")
	flag.Parse()
	if len(os.Args) < 2 {
		log.Fatal("Incorrect arguments! Use --help to get the usage.")
	}
	cut = !scale
	args := flag.Args()

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
		fmt.Println("output ios launch images")
		for _, root := range args {
			GenerateLaunchImage(root)
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
