package main

import (
	"errors"
	"fmt"
	"github.com/go-fsnotify/fsnotify"
	"github.com/nfnt/resize"
	flag "github.com/ogier/pflag"
	"github.com/oliamb/cutter"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	MaxWidth  int = 620
	MaxHeight int = 960
)

var (
	watcher        *fsnotify.Watcher
	cut            bool
	scale          bool
	appicon        bool
	launchimage    bool
	watch          bool
	imagedirectory string
)

func init() {
	flag.BoolVar(&cut, "cut", false, "cut mode")
	flag.BoolVar(&cut, "c", false, "cut mode")
	flag.BoolVar(&scale, "scale", false, "scale mode")
	flag.BoolVar(&scale, "s", false, "scale mode")
	flag.BoolVar(&appicon, "appicon", false, "generate app icon")
	flag.BoolVar(&appicon, "a", false, "generate app icon")
	flag.BoolVar(&launchimage, "launchimage", false, "generate launch images")
	flag.BoolVar(&launchimage, "l", false, "generate launch images")
	flag.BoolVar(&watch, "watch", false, "watch directories change")
	flag.BoolVar(&watch, "w", false, "watch directories change")
}

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

func traverseCut(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		doCutImage(path)
	} else {
		fmt.Printf("Visited: %s\n", path)
	}
	return nil
}

func doCutImage(path string) {
	if strings.LastIndex(path, ")-p.jpg") < 0 &&
		strings.LastIndex(path, ")-p.png") < 0 &&
		strings.LastIndex(path, "-m.jpg") < 0 &&
		strings.LastIndex(path, "-m.png") < 0 {
		savePath := path + "(1)-p.jpg"
		if _, err := os.Stat(savePath); err != nil {
			// not exists
			cutImage(path, 2)
		}

		savePath = path + "(1)-p.png"
		if _, err := os.Stat(savePath); err != nil {
			// not exists
			cutImage(path, 1)
		}
	}
}

func cutImage(filepath string, imagetype int) error {
	reader, err := os.Open(filepath)
	if err != nil {
		log.Println(filepath, err)
		return err
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Println(filepath, err)
		return err
	}
	bounds := m.Bounds()

	expectHeight := bounds.Size().X * MaxHeight / MaxWidth
	if bounds.Size().Y > expectHeight {
		blockCount := bounds.Size().Y/expectHeight + 1

		for i := 0; i < blockCount && i*expectHeight < bounds.Size().Y; i++ {
			// cut first
			croppedImg, err := cutter.Crop(m, cutter.Config{
				Width:  bounds.Size().X,
				Height: expectHeight,
				Anchor: image.Point{0, i * expectHeight},
			})

			if err != nil {
				log.Println(filepath, err)
				return err
			}

			// resize then
			im := resize.Resize(uint(MaxWidth), 0, croppedImg, resize.Bilinear)

			// save to file finally
			var savePath string
			switch imagetype {
			case 1:
				savePath = fmt.Sprintf("%s(%d)-p.png", filepath, i+1)
			default:
				savePath = fmt.Sprintf("%s(%d)-p.jpg", filepath, i+1)
			}
			fmt.Printf("%s width=%d, height=%d, cropped to %s\n",
				filepath, bounds.Size().X, bounds.Size().Y, savePath)

			if err := saveImage(&im, savePath, imagetype); err != nil {
				return err
			}
		}
	}
	return nil
}

func scaleImage(filepath string, imagetype int) error {
	reader, err := os.Open(filepath)
	if err != nil {
		log.Println(filepath, err)
		return err
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Println(filepath, err)
		return err
	}
	bounds := m.Bounds()

	if bounds.Size().X > MaxWidth || bounds.Size().Y > MaxHeight {
		// scale it
		var im image.Image
		if bounds.Size().X > MaxWidth &&
			bounds.Size().Y*MaxWidth/bounds.Size().X < MaxHeight {
			im = resize.Resize(uint(MaxWidth), 0, m, resize.Bilinear)
		} else {
			im = resize.Resize(0, uint(MaxHeight), m, resize.Bilinear)
		}
		var savePath string
		switch imagetype {
		case 1:
			savePath = filepath + "-m.png"
		default:
			savePath = filepath + "-m.jpg"
		}

		fmt.Printf("%s width=%d, height=%d, scaled to %s\n",
			filepath, bounds.Size().X, bounds.Size().Y, savePath)

		return saveImage(&im, savePath, imagetype)
	}

	return nil
}

func traverseScale(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		doScaleImage(path)
	} else {
		fmt.Printf("Visited: %s\n", path)
	}
	return nil
}

func doScaleImage(path string) {
	if strings.LastIndex(path, "-m.jpg") < 0 &&
		strings.LastIndex(path, "-m.png") < 0 {
		savePath := path + "-m.jpg"
		if _, err := os.Stat(savePath); err != nil {
			// not exists
			scaleImage(path, 2)
		}

		savePath = path + "-m.png"
		if _, err := os.Stat(savePath); err != nil {
			// not exists
			scaleImage(path, 1)
		}
	}
}

func watchDirectory(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		if err := watcher.Add(path); err != nil {
			log.Println(err)
		}
	} else {
		doScaleImage(path)
	}
	return nil
}

func main() {
	flag.Parse()
	if len(os.Args) < 2 {
		log.Fatal(errors.New("Incorrect arguments!"))
	}
	args := flag.Args()

	if watch {
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
	} else if cut {
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
	} else if appicon {
		fmt.Println("output ios app icons")
	} else if launchimage {
		fmt.Println("output ios launch images")
	} else {
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
}
