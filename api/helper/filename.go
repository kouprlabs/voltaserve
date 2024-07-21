package helper

import (
	"fmt"
	"path/filepath"
)

func UniqueFilename(name string) string {
	return fmt.Sprintf("%s %s%s", FilenameWithoutExtension(name), NewID(), filepath.Ext(name))
}

func FilenameWithoutExtension(name string) string {
	withExt := filepath.Base(name)
	return withExt[0 : len(withExt)-len(filepath.Ext(name))]
}
