package identifier

import (
	"path/filepath"
	"strings"
)

type ImageIdentifier struct{}

func NewImageIdentifier() *ImageIdentifier {
	return &ImageIdentifier{}
}

func (ii ImageIdentifier) IsJPEG(path string) bool {
	path = strings.ToLower(path)
	return filepath.Ext(path) == ".jpg" ||
		filepath.Ext(path) == ".jpeg" ||
		filepath.Ext(path) == ".jpe" ||
		filepath.Ext(path) == ".jfif" ||
		filepath.Ext(path) == ".jif"
}

func (ii ImageIdentifier) IsPNG(path string) bool {
	path = strings.ToLower(path)
	return filepath.Ext(path) == ".png"
}

func (ii ImageIdentifier) IsTIFF(path string) bool {
	path = strings.ToLower(path)
	return filepath.Ext(path) == ".tiff" ||
		filepath.Ext(path) == ".tif"
}
