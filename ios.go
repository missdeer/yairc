package main

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/missdeer/golib/fsutil"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

type handler func(image.Image, image.Image, string, *launchImageSpec) error

type launchImageSpec struct {
	Width   int
	Height  int
	Postfix string
	Handler handler
}

type appIconSpec struct {
	Length int
	Name   string
}

var (
	launchImageSpecifications = []launchImageSpec{
		{640, 960, "LaunchImage-iOS7@2x~iphone.png", BackgroundForegroundHandler},
		{960, 640, "LaunchImage-iOS7-Landscape@2x~iphone.png", BackgroundForegroundHandler},
		{640, 1136, "LaunchImage-iOS7-568h@2x~iphone.png", BackgroundForegroundHandler},
		{1136, 640, "LaunchImage-iOS7-Landscape-568h@2x~iphone.png", BackgroundForegroundHandler},
		{750, 1334, "LaunchImage-375w-667h@2x~iphone.png", BackgroundForegroundHandler},
		{1334, 750, "LaunchImage-Landscape-375w-667h@2x~iphone.png", BackgroundForegroundHandler},
		{1242, 2208, "LaunchImage-414w-736h@3x~iphone.png", BackgroundForegroundHandler},
		{2208, 1242, "LaunchImage-Landscape-414w-736h@3x~iphone.png", BackgroundForegroundHandler},
		{768, 1024, "LaunchImage-iOS7-Portrait~ipad.png", BackgroundForegroundHandler},
		{1024, 768, "LaunchImage-iOS7-Landscape~ipad.png", BackgroundForegroundHandler},
		{1536, 2048, "LaunchImage-iOS7-Portrait@2x~ipad.png", BackgroundForegroundHandler},
		{2048, 1536, "LaunchImage-iOS7-Landscape@2x~ipad.png", BackgroundForegroundHandler},
		{1668, 2224, "LaunchImage-Portrait-1112@2x.png", BackgroundForegroundHandler},
		{2224, 1668, "LaunchImage-Landscape-1112@2x.png", BackgroundForegroundHandler},
		{1125, 2436, "LaunchImage-375w-812h@3x.png", BackgroundForegroundHandler},
		{2436, 1125, "LaunchImage-Landscape-375w-812h@3x.png", BackgroundForegroundHandler},
		{2048, 2732, "LaunchImage-Portrait@2x.png", BackgroundForegroundHandler},
		{2732, 2048, "LaunchImage-Landscape@2x.png", BackgroundForegroundHandler},
	}
	appIconSpecifications = []appIconSpec{
		// Recommended if you have a Settings bundle, optional otherwise
		{29, "AppIcon29x29.png"},
		{29, "AppIcon29x29.png"},
		{58, "AppIcon29x29@2x.png"},
		{58, "AppIcon29x29@2x.png"},
		{87, "AppIcon29x29@3x.png"},
		{57, "AppIcon57x57.png"},
		{114, "AppIcon57x57@2x.png"},
		{144, "AppIcon72x72@2x.png"},
		{72, "AppIcon72x72.png"},
		// Spotlight
		{40, "AppIcon40x40.png"},
		{80, "AppIcon40x40@2x.png"},
		{80, "AppIcon40x40@2x.png"},
		{50, "AppIcon50x50.png"},
		{100, "AppIcon50x50@2x.png"},
		{120, "AppIcon40x40@3x.png"},
		// Home screen on iPad
		{76, "AppIcon76x76.png"},
		{152, "AppIcon76x76@2x.png"},
		// Home screen on iPad Pro
		{167, "AppIcon83.5x83.5@2x.png"},
		// Home screen on iPhone/iPod Touch with retina display
		{20, "AppIcon20x20.png"},
		{40, "AppIcon20x20@2x.png"},
		{60, "AppIcon20x20@3x.png"},
		{120, "AppIcon60x60@2x.png"},
		{180, "AppIcon60x60@3x.png"},
		// iWatch
		{48, "AppIcon24@2x.png"},
		{55, "AppIcon27.5@2x.png"},
		{58, "AppIcon29@2x.png"},
		{80, "AppIcon40@2x.png"},
		{87, "AppIcon29@3x.png"},
		{88, "AppIcon44@2x.png"},
		{172, "AppIcon86@2x.png"},
		{196, "AppIcon98@2x.png"},
		// App list in iTunes
		// {512, "iTunesArtwork.png"},
		{1024, "iTunesArtwork@2x.png"},
	}

	contentsJson = `{
	"images": [
		{
			"size": "20x20",
			"idiom": "iphone",
			"filename": "AppIcon20x20@2x.png",
			"scale": "2x"
		},
		{
			"size": "20x20",
			"idiom": "iphone",
			"filename": "AppIcon20x20@3x.png",
			"scale": "3x"
		},
		{
			"size": "29x29",
			"idiom": "iphone",
			"filename": "AppIcon29x29@2x.png",
			"scale": "2x"
		},
		{
			"size": "29x29",
			"idiom": "iphone",
			"filename": "AppIcon29x29@3x.png",
			"scale": "3x"
		},
		{
			"size": "40x40",
			"idiom": "iphone",
			"filename": "AppIcon40x40@2x.png",
			"scale": "2x"
		},
		{
			"size": "40x40",
			"idiom": "iphone",
			"filename": "AppIcon40x40@3x.png",
			"scale": "3x"
		},
		{
			"size": "60x60",
			"idiom": "iphone",
			"filename": "AppIcon60x60@2x.png",
			"scale": "2x"
		},
		{
			"size": "60x60",
			"idiom": "iphone",
			"filename": "AppIcon60x60@3x.png",
			"scale": "3x"
		},
		{
			"size": "20x20",
			"idiom": "ipad",
			"filename": "AppIcon20x20.png",
			"scale": "1x"
		},
		{
			"size": "20x20",
			"idiom": "ipad",
			"filename": "AppIcon20x20@2x.png",
			"scale": "2x"
		},
		{
			"size": "29x29",
			"idiom": "ipad",
			"filename": "AppIcon29x29.png",
			"scale": "1x"
		},
		{
			"size": "29x29",
			"idiom": "ipad",
			"filename": "AppIcon29x29@2x.png",
			"scale": "2x"
		},
		{
			"size": "40x40",
			"idiom": "ipad",
			"filename": "AppIcon40x40.png",
			"scale": "1x"
		},
		{
			"size": "40x40",
			"idiom": "ipad",
			"filename": "AppIcon40x40@2x.png",
			"scale": "2x"
		},
		{
			"size": "76x76",
			"idiom": "ipad",
			"filename": "AppIcon76x76.png",
			"scale": "1x"
		},
		{
			"size": "76x76",
			"idiom": "ipad",
			"filename": "AppIcon76x76@2x.png",
			"scale": "2x"
		},
		{
			"size": "83.5x83.5",
			"idiom": "ipad",
			"filename": "AppIcon83.5x83.5@2x.png",
			"scale": "2x"
		},
		{
			"size": "1024x1024",
			"idiom": "ios-marketing",
			"filename": "iTunesArtwork@2x.png",
			"scale": "1x"
		},
		{
			"size": "24x24",
			"idiom": "watch",
			"scale": "2x",
			"filename": "AppIcon24@2x.png",
			"role": "notificationCenter",
			"subtype": "38mm"
		},
		{
			"size": "27.5x27.5",
			"idiom": "watch",
			"scale": "2x",
			"filename": "AppIcon27.5@2x.png",
			"role": "notificationCenter",
			"subtype": "42mm"
		},
		{
			"size": "29x29",
			"idiom": "watch",
			"filename": "AppIcon29@2x.png",
			"role": "companionSettings",
			"scale": "2x"
		},
		{
			"size": "29x29",
			"idiom": "watch",
			"filename": "AppIcon29@3x.png",
			"role": "companionSettings",
			"scale": "3x"
		},
		{
			"size": "40x40",
			"idiom": "watch",
			"scale": "2x",
			"filename": "AppIcon40@2x.png",
			"role": "appLauncher",
			"subtype": "38mm"
		},
		{
			"size": "44x44",
			"idiom": "watch",
			"scale": "2x",
			"filename": "AppIcon44@2x.png",
			"role": "longLook",
			"subtype": "42mm"
		},
		{
			"size": "86x86",
			"idiom": "watch",
			"scale": "2x",
			"filename": "AppIcon86@2x.png",
			"role": "quickLook",
			"subtype": "38mm"
		},
		{
			"size": "98x98",
			"idiom": "watch",
			"scale": "2x",
			"filename": "AppIcon98@2x.png",
			"role": "quickLook",
			"subtype": "42mm"
		}
	],
	"properties": {
		"pre-rendered": true
	}
}`
)

func BackgroundForegroundHandler(bm image.Image, fm image.Image, savePath string, spec *launchImageSpec) error {
	im := resize.Resize(0, uint(spec.Height), bm, resize.Bilinear)
	if im.Bounds().Size().X < spec.Width {
		im = resize.Resize(uint(spec.Width), 0, bm, resize.Bilinear)
		log.Println("resize by width")
	} else {
		log.Println("resize by height")
	}

	var err error
	if im.Bounds().Size().X > spec.Width {
		im, err = cutter.Crop(im, cutter.Config{
			Width:  spec.Width,
			Height: spec.Height,
			Anchor: image.Point{(im.Bounds().Size().X - spec.Width) / 2, 0},
		})
		if err != nil {
			log.Println(savePath, err)
		}
	}
	if im.Bounds().Size().Y > spec.Height {
		im, err = cutter.Crop(im, cutter.Config{
			Width:  spec.Width,
			Height: spec.Height,
			Anchor: image.Point{0, (im.Bounds().Size().Y - spec.Height) / 2},
		})
		if err != nil {
			log.Println(savePath, err)
		}
	}

	m := image.NewRGBA(image.Rect(0, 0, spec.Width, spec.Height))
	draw.Draw(m, m.Bounds(), im, im.Bounds().Min, draw.Src)
	if spec.Width < spec.Height {
		x := spec.Width / 4
		y := spec.Height/2 - x
		sm := resize.Resize(uint(x*2), 0, fm, resize.Bilinear)
		draw.Draw(m, image.Rect(x, y, x*3, y+x*2), sm, image.Point{0, 0}, draw.Over)
	} else {
		y := spec.Height / 4
		x := spec.Width/2 - y
		sm := resize.Resize(0, uint(y*2), fm, resize.Bilinear)
		draw.Draw(m, image.Rect(x, y, x+y*2, y*3), sm, image.Point{0, 0}, draw.Over)
	}

	if err = saveRGBA(m, savePath, 1); err != nil {
		log.Println(savePath, err)
	}
	return nil
}

func GenerateLaunchImage() error {
	reader, err := os.Open(backgroundImagePath)
	if err != nil {
		log.Println(backgroundImagePath, err)
		return err
	}
	defer reader.Close()
	bm, _, err := image.Decode(reader)
	if err != nil {
		log.Println(backgroundImagePath, err)
		return err
	}

	reader, err = os.Open(foregroundImagePath)
	if err != nil {
		log.Println(foregroundImagePath, err)
		return err
	}
	defer reader.Close()
	fm, _, err := image.Decode(reader)
	if err != nil {
		log.Println(foregroundImagePath, err)
		return err
	}

	os.Mkdir("launchimage", 0755)
	os.Mkdir("launchimage/ios", 0755)
	for _, spec := range launchImageSpecifications {
		savePath := "launchimage/ios/" + spec.Postfix

		log.Println("generating ", savePath)
		spec.Handler(bm, fm, savePath, &spec)
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

	origLength := m.Bounds().Dx()
	length := origLength * 4 / 5
	m = resize.Resize(uint(length), uint(length), m, resize.Bilinear)

	bm := image.NewRGBA(image.Rect(0, 0, origLength, origLength))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(bm, bm.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)
	draw.Draw(bm, image.Rect(origLength/10, origLength/10, origLength/10+length, origLength/10+length), m, image.Point{0, 0}, draw.Over)

	os.Mkdir("appicon", 0755)
	os.Mkdir("appicon/ios", 0755)
	os.Mkdir("appicon/ios/Images.xcassets", 0755)
	os.Mkdir("appicon/ios/Images.xcassets/AppIcon.appiconset", 0755)
	for _, spec := range appIconSpecifications {
		im := resize.Resize(uint(spec.Length), uint(spec.Length), bm, resize.Bilinear)
		if err := saveImage(&im, "appicon/ios/Images.xcassets/AppIcon.appiconset/"+spec.Name, 1); err != nil {
			log.Println(spec.Name, err)
		}
	}

	fd, err := os.OpenFile("appicon/ios/Images.xcassets/AppIcon.appiconset/Contents.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	fd.WriteString(contentsJson)
	fd.Close()
	return nil
}

func iconScale(inputFile string, outputDir string) error {
	if b, e := fsutil.DirExists(outputDir); e != nil || !b {
		if e = os.MkdirAll(outputDir, 0755); e != nil {
			log.Fatal(e)
			return e
		}
	}
	if b, e := fsutil.DirExists(path.Join(outputDir, "x18")); e != nil || !b {
		if e = os.MkdirAll(path.Join(outputDir, "x18"), 0755); e != nil {
			log.Fatal(e)
			return e
		}
	}
	if b, e := fsutil.DirExists(path.Join(outputDir, "x36")); e != nil || !b {
		if e = os.MkdirAll(path.Join(outputDir, "x36"), 0755); e != nil {
			log.Fatal(e)
			return e
		}
	}
	if b, e := fsutil.DirExists(path.Join(outputDir, "x48")); e != nil || !b {
		if e = os.MkdirAll(path.Join(outputDir, "x48"), 0755); e != nil {
			log.Fatal(e)
			return e
		}
	}
	reader, err := os.Open(inputFile)
	if err != nil {
		log.Println(err)
		return err
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Println(err)
		return err
	}
	infos := []struct {
		length        uint
		relativePaths []string
	}{
		{18, []string{filepath.Join(outputDir, "x18", filepath.Base(inputFile))}},
		{24, []string{filepath.Join(outputDir, filepath.Base(inputFile))}},
		{36, []string{
			filepath.Join(outputDir, "x18", filepath.Base(inputFile)[:len(filepath.Base(inputFile))-len(filepath.Ext(inputFile))]+"@2x"+filepath.Ext(inputFile)),
			filepath.Join(outputDir, "x36", filepath.Base(inputFile))}},
		{48, []string{
			filepath.Join(outputDir, filepath.Base(inputFile)[:len(filepath.Base(inputFile))-len(filepath.Ext(inputFile))]+"@2x"+filepath.Ext(inputFile)),
			filepath.Join(outputDir, "x48", filepath.Base(inputFile))}},
		{54, []string{filepath.Join(outputDir, "x18", filepath.Base(inputFile)[:len(filepath.Base(inputFile))-len(filepath.Ext(inputFile))]+"@3x"+filepath.Ext(inputFile))}},
		{72, []string{
			filepath.Join(outputDir, filepath.Base(inputFile)[:len(filepath.Base(inputFile))-len(filepath.Ext(inputFile))]+"@3x"+filepath.Ext(inputFile)),
			filepath.Join(outputDir, "x18", filepath.Base(inputFile)[:len(filepath.Base(inputFile))-len(filepath.Ext(inputFile))]+"@4x"+filepath.Ext(inputFile)),
			filepath.Join(outputDir, "x36", filepath.Base(inputFile)[:len(filepath.Base(inputFile))-len(filepath.Ext(inputFile))]+"@2x"+filepath.Ext(inputFile))}},
		{96, []string{
			filepath.Join(outputDir, filepath.Base(inputFile)[:len(filepath.Base(inputFile))-len(filepath.Ext(inputFile))]+"@4x"+filepath.Ext(inputFile)),
			filepath.Join(outputDir, "x48", filepath.Base(inputFile)[:len(filepath.Base(inputFile))-len(filepath.Ext(inputFile))]+"@2x"+filepath.Ext(inputFile))}},
		{108, []string{filepath.Join(outputDir, "x36", filepath.Base(inputFile)[:len(filepath.Base(inputFile))-len(filepath.Ext(inputFile))]+"@3x"+filepath.Ext(inputFile))}},
		{144, []string{
			filepath.Join(outputDir, "x36", filepath.Base(inputFile)[:len(filepath.Base(inputFile))-len(filepath.Ext(inputFile))]+"@4x"+filepath.Ext(inputFile)),
			filepath.Join(outputDir, "x48", filepath.Base(inputFile)[:len(filepath.Base(inputFile))-len(filepath.Ext(inputFile))]+"@3x"+filepath.Ext(inputFile))}},
		{192, []string{filepath.Join(outputDir, "x48", filepath.Base(inputFile)[:len(filepath.Base(inputFile))-len(filepath.Ext(inputFile))]+"@4x"+filepath.Ext(inputFile))}},
	}
	for _, info := range infos {
		im := resize.Resize(info.length, info.length, m, resize.Bilinear)
		for _, relativePath := range info.relativePaths {
			if err := saveImage(&im, relativePath, 1); err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}

func iOSScale(origin string, templateSize string) error {
	switch templateSize {
	case "1x":
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

		name := filepath.Join(filepath.Dir(origin), filepath.Base(origin)[:len(filepath.Base(origin))-len(filepath.Ext(origin))]+"@2x"+filepath.Ext(origin))
		im := resize.Resize(uint(m.Bounds().Size().X*2), uint(m.Bounds().Size().Y*2), m, resize.Bilinear)
		if err := saveImage(&im, name, 1); err != nil {
			log.Println(name, err)
		}

		name = filepath.Join(filepath.Dir(origin), filepath.Base(origin)[:len(filepath.Base(origin))-len(filepath.Ext(origin))]+"@3x"+filepath.Ext(origin))
		im = resize.Resize(uint(m.Bounds().Size().X*3), uint(m.Bounds().Size().Y*3), m, resize.Bilinear)
		if err := saveImage(&im, name, 1); err != nil {
			log.Println(name, err)
		}
	case "2x":
		var one, two, three string
		if !strings.HasSuffix(filepath.Base(origin), "@2x") {
			one = origin
			two = filepath.Join(filepath.Dir(origin), filepath.Base(origin)[:len(filepath.Base(origin))-len(filepath.Ext(origin))]+"@2x"+filepath.Ext(origin))
			three = filepath.Join(filepath.Dir(origin), filepath.Base(origin)[:len(filepath.Base(origin))-len(filepath.Ext(origin))]+"@3x"+filepath.Ext(origin))
			os.Rename(origin, two)
		} else {
			one = filepath.Join(filepath.Dir(origin), filepath.Base(origin)[:len(filepath.Base(origin))-len(filepath.Ext(origin))-3]+filepath.Ext(origin))
			two = origin
			three = filepath.Join(filepath.Dir(origin), filepath.Base(origin)[:len(filepath.Base(origin))-len(filepath.Ext(origin))-3]+"@3x"+filepath.Ext(origin))
		}

		reader, err := os.Open(two)
		if err != nil {
			log.Println(two, err)
			return err
		}
		defer reader.Close()
		m, _, err := image.Decode(reader)
		if err != nil {
			log.Println(two, err)
			return err
		}

		im := resize.Resize(uint(m.Bounds().Size().X/2), uint(m.Bounds().Size().Y/2), m, resize.Bilinear)
		if err := saveImage(&im, one, 1); err != nil {
			log.Println(one, err)
		}

		im = resize.Resize(uint(m.Bounds().Size().X*3/2), uint(m.Bounds().Size().Y*3/2), m, resize.Bilinear)
		if err := saveImage(&im, three, 1); err != nil {
			log.Println(three, err)
		}
	case "3x":
		var one, two, three string
		if !strings.HasSuffix(filepath.Base(origin), "@3x") {
			one = origin
			two = filepath.Join(filepath.Dir(origin), filepath.Base(origin)[:len(filepath.Base(origin))-len(filepath.Ext(origin))]+"@2x"+filepath.Ext(origin))
			three = filepath.Join(filepath.Dir(origin), filepath.Base(origin)[:len(filepath.Base(origin))-len(filepath.Ext(origin))]+"@3x"+filepath.Ext(origin))
			os.Rename(origin, three)
		} else {
			one = filepath.Join(filepath.Dir(origin), filepath.Base(origin)[:len(filepath.Base(origin))-len(filepath.Ext(origin))-3]+filepath.Ext(origin))
			two = filepath.Join(filepath.Dir(origin), filepath.Base(origin)[:len(filepath.Base(origin))-len(filepath.Ext(origin))-3]+"@2x"+filepath.Ext(origin))
			three = origin
		}

		reader, err := os.Open(three)
		if err != nil {
			log.Println(three, err)
			return err
		}
		defer reader.Close()
		m, _, err := image.Decode(reader)
		if err != nil {
			log.Println(three, err)
			return err
		}

		im := resize.Resize(uint(m.Bounds().Size().X/3), uint(m.Bounds().Size().Y/3), m, resize.Bilinear)
		if err := saveImage(&im, one, 1); err != nil {
			log.Println(one, err)
		}

		im = resize.Resize(uint(m.Bounds().Size().X*2/3), uint(m.Bounds().Size().Y*2/3), m, resize.Bilinear)
		if err := saveImage(&im, two, 1); err != nil {
			log.Println(two, err)
		}
	default:
		log.Fatal("unrecognized template size")
	}

	return nil
}
