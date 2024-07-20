package helper

import (
	"fmt"
	"path/filepath"
)

func UniqueFilename(name string) string {
	return fmt.Sprintf("%s %s%s", filepath.Base(name), NewID(), filepath.Ext(name))
}
