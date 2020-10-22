package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/missdeer/yairc/util"
	flag "github.com/spf13/pflag"
)

var (
	compress               bool
	backgroundImagePath    string
	foregroundImagePath    string
	inputPath              string
	outputPath             string
	action                 string
	platform                      = "ios"
	red                    uint32 = 127
	green                  uint32 = 127
	blue                   uint32 = 127
	outputHeight           uint
	outputWidth            uint
	transparentWhiteDirect bool
	// Gitcommit contains the commit where we built from.
	GitCommit string

	imageFormatMap = map[string]int{
		".icns": util.IT_icns,
		".png":  util.IT_png,
		".jpg":  util.IT_jpeg,
		".jpeg": util.IT_jpeg,
		".gif":  util.IT_gif,
		".webp": util.IT_webp,
		".tiff": util.IT_tiff,
	}
)

func main() {
	showHelpMessage := false
	showVersion := false
	flag.Uint32VarP(&red, "red", "", red, "set red threshold")
	flag.Uint32VarP(&green, "green", "", green, "set green threshold")
	flag.Uint32VarP(&blue, "blue", "", blue, "set blue threshold")
	flag.BoolVarP(&compress, "compress", "", true, "compress output PNG files")
	flag.StringVarP(&platform, "platform", "p", "common", "candidates: ios, android, common")
	flag.StringVarP(&action, "action", "a", "", "candidats: icons, appIcon, launchImage, transparent, invert, resize, convert, info")
	flag.StringVarP(&backgroundImagePath, "background", "b", "", "path of background image for launch image")
	flag.StringVarP(&foregroundImagePath, "foreground", "f", "", "path of foreground image for launch image")
	flag.StringVarP(&inputPath, "input", "i", "", "input image file path")
	flag.StringVarP(&outputPath, "output", "o", ".", "output directory/file path")
	flag.UintVarP(&outputHeight, "height", "", 0, "set output image height, 0 for original height")
	flag.UintVarP(&outputWidth, "width", "", 0, "set output image width, 0 for original width")
	flag.BoolVarP(&transparentWhiteDirect, "transparent-white-direct", "", false, "false - make white color be transparent, true - make black color be transparent")
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
	if action == "icons" && inputPath != "" && outputPath != "" {
		log.Println("generate /@2x/@3x/@4x & /x18/x36/x48 icons from", inputPath, "to", outputPath)
		err := iconScale(inputPath, outputPath)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// ios app icon mode
	if action == "appIcon" && platform == "ios" {
		fmt.Println("output ios app icons")
		err := GenerateAppIcon(inputPath)
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
		err := GenerateLauncherIcon(inputPath)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if action == "launchImage" && platform == "ios" {
		fmt.Println("output android splash screen images")
		err := GenerateSplashScreen(inputPath)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// convert to .icns file
	if action == "convert" && inputPath != "" {
		dir := filepath.Dir(outputPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, 0755); err != nil {
				log.Fatal("mkdir failed", err)
			}
		}
		dest, err := os.Create(outputPath)
		if err != nil {
			log.Fatal("opening destination file failed", err)
		}
		defer dest.Close()
		inFile, err := util.OpenURI(inputPath)
		if err != nil {
			log.Fatal("opening source image failed", err)
		}
		defer inFile.Close()
		srcImg, _, err := util.ImageDecode(inFile)
		if err != nil {
			log.Fatal("decoding source image failed", err)
		}
		outExt := filepath.Ext(outputPath)
		if it, ok := imageFormatMap[strings.ToLower(outExt)]; ok {
			err = util.SaveImage(srcImg, outputPath, it)
		} else {
			log.Fatal("unsupported target image format")
		}
		if err != nil {
			log.Fatal("encoding failed", err)
		}
		return
	}

	args := flag.Args()
	if inputPath != "" {
		args = append(args, inputPath)
	}

	if action == "transparent" && len(args) > 0 {
		log.Println("transparent color")
		for _, uri := range args {
			im, err := util.Transparent(uri, red, green, blue, transparentWhiteDirect)
			if err != nil {
				log.Println(err)
				continue
			}
			fn := outputPath
			if fn == "" {
				fn = inputPath[:len(inputPath)-len(filepath.Ext(inputPath))] + ".transparent.png"
			}
			outExt := filepath.Ext(fn)
			if it, ok := imageFormatMap[strings.ToLower(outExt)]; ok {
				if err = util.SaveImage(im, outputPath, it); err != nil {
					log.Println("encoding failed", err)
					continue
				}

				if it == util.IT_png {
					if err = util.DoCrush(compress, fn); err != nil {
						log.Println(err)
					}
				}
			} else {
				log.Println("unsupported target image format")
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
			fn := outputPath
			if fn == "" {
				fn = inputPath[:len(inputPath)-len(filepath.Ext(inputPath))] + ".invert.png"
			}
			outExt := filepath.Ext(fn)
			if it, ok := imageFormatMap[strings.ToLower(outExt)]; ok {
				if err = util.SaveImage(im, outputPath, it); err != nil {
					log.Println("encoding failed", err)
					continue
				}

				if it == util.IT_png {
					if err = util.DoCrush(compress, fn); err != nil {
						log.Println(err)
					}
				}
			} else {
				log.Println("unsupported target image format")
			}
		}
		return
	}

	if action == "resize" && len(args) > 0 {
		log.Println("resize images")
		for _, uri := range args {
			im, err := util.Resize(uri, outputWidth, outputHeight)
			if err != nil {
				log.Println(err)
				continue
			}
			fn := outputPath
			if fn == "" {
				fn = inputPath[:len(inputPath)-len(filepath.Ext(inputPath))] + ".resized.png"
			}
			outExt := filepath.Ext(fn)
			if it, ok := imageFormatMap[strings.ToLower(outExt)]; ok {
				if err = util.SaveImage(im, outputPath, it); err != nil {
					log.Println("encoding failed", err)
					continue
				}

				if it == util.IT_png {
					if err = util.DoCrush(compress, fn); err != nil {
						log.Println(err)
					}
				}
			} else {
				log.Println("unsupported target image format")
			}
		}
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
				r, g, b, a := k.RGBA()
				fmt.Printf("r=%3d,g=%3d,b=%3d,a=%3d: count:=%d\n", r, g, b, a, v)
			}
			fmt.Println("===============================")
		}
		return
	}
}
