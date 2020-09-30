package util

import "os"

func IsDir(path string) (bool, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return false, err
	}
	file := f

	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}
