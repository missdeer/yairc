package main

import (
	"errors"
	"fmt"
	"github.com/howeyc/fsnotify"
	"github.com/nfnt/resize"
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
	watcher *fsnotify.Watcher
)

func scale(filepath string, imagetype int) error {
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
		if bounds.Size().X > MaxWidth && bounds.Size().Y*MaxWidth/bounds.Size().X < MaxHeight {
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

		fmt.Printf("%s width=%d, height=%d, saved to %s \n", filepath, bounds.Size().X, bounds.Size().Y, savePath)

		var file *os.File
		if f, err := os.OpenFile(savePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
			log.Fatal(savePath, err)
		} else {
			file = f
		}
		defer file.Close()

		switch imagetype {
		case 1:
			err = png.Encode(file, im)
		default:
			err = jpeg.Encode(file, im, &jpeg.Options{100})
		}

		if err != nil {
			log.Println(savePath, err)
			return err
		}
	}

	return nil
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		do_scale(path)
	} else {
		fmt.Printf("Visited: %s\n", path)
	}
	return nil
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

func do_scale(path string) {
	if strings.LastIndex(path, "-m.jpg") < 0 && strings.LastIndex(path, "-m.png") < 0 {
		savePath := path + "-m.jpg"
		if _, err := os.Stat(savePath); err != nil {
			// not exists
			scale(path, 2)
		}

		savePath = path + "-m.png"
		if _, err := os.Stat(savePath); err != nil {
			// not exists
			scale(path, 1)
		}
	}
}

func watch_directory(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		if err := watcher.WatchFlags(path, fsnotify.FSN_DELETE|fsnotify.FSN_MODIFY); err != nil {
			log.Println(err)
		}
	} else {
		do_scale(path)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("Incorrect arguments!"))
	}

	if os.Args[1] == `-w` || os.Args[1] == `--watch` {
		var err error
		watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(errors.New("creating filesystem watcher failed"))
		}

		go func() {
			for {
				select {
				case event := <-watcher.Event:
					if b, e := isDir(event.Name); e == nil && b == false {
						if event.IsDelete() || event.IsModify() {
							// delete associated files
							if strings.LastIndex(event.Name, "-m.jpg") < 0 && strings.LastIndex(event.Name, "-m.png") < 0 {
								os.Remove(event.Name + "-m.jpg")
								os.Remove(event.Name + "-m.png")
							}
						}
						if event.IsModify() {
							do_scale(event.Name)
						}
					}
				case err := <-watcher.Error:
					log.Println("error:", err)
				}
			}
		}()

		for i := 2; i < len(os.Args); i++ {
			root := os.Args[i]
			if b, e := isDir(root); e == nil && b == true {
				err := filepath.Walk(root, watch_directory)
				if err != nil {
					log.Println(err)
					continue
				}
			} else {
				do_scale(root)
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
	} else {
		for i := 1; i < len(os.Args); i++ {
			root := os.Args[i]
			if b, e := isDir(root); e == nil && b == true {
				err := filepath.Walk(root, visit)
				if err != nil {
					log.Println(err)
					continue
				}
			} else {
				do_scale(root)
			}
		}
	}
}
