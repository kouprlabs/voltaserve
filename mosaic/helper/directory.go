package helper

import (
	"os"
	"path/filepath"
)

func FindFileWithExtension(dirPath string, ext string) (*string, error) {
	var res string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ext {
			res = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if res == "" {
		return nil, nil
	}
	return &res, nil
}
