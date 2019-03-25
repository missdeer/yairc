package main

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	"github.com/ultimate-guitar/go-imagequant"
)

func crush(filePath string) error {
	suffix := uuid.New().String()
	err := crushFile(filePath, filePath+"."+suffix, 3, png.BestCompression)
	if err != nil {
		return err
	}
	err = os.Remove(filePath)
	if err != nil {
		return err
	}
	err = os.Rename(filePath+"."+suffix, filePath)
	if err != nil {
		return err
	}

	return nil
}

func crushFile(sourcefile, destfile string, speed int, compression png.CompressionLevel) error {

	sourceFh, err := os.OpenFile(sourcefile, os.O_RDONLY, 0444)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %s", err.Error())
	}
	defer sourceFh.Close()

	image, err := ioutil.ReadAll(sourceFh)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}

	optiImage, err := imagequant.Crush(image, speed, compression)
	if err != nil {
		return fmt.Errorf("imagequant.Crush: %s", err.Error())
	}

	destFh, err := os.OpenFile(destfile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %s", err.Error())
	}
	defer destFh.Close()

	destFh.Write(optiImage)
	return nil
}
