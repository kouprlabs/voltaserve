package identifier

import (
	"path/filepath"
	"strings"
)

type FileIdentifier struct {
}

func NewFileIdentifier() *FileIdentifier {
	return &FileIdentifier{}
}

func (fi *FileIdentifier) IsPDF(path string) bool {
	return strings.ToLower(filepath.Ext(path)) == ".pdf"
}

func (fi *FileIdentifier) IsOffice(path string) bool {
	extensions := []string{
		".xls",
		".doc",
		".ppt",
		".xlsx",
		".docx",
		".pptx",
		".odt",
		".ott",
		".ods",
		".ots",
		".odp",
		".otp",
		".odg",
		".otg",
		".odf",
		".odc",
		".rtf",
	}
	extension := filepath.Ext(path)
	for _, v := range extensions {
		if strings.ToLower(extension) == v {
			return true
		}
	}
	return false
}

func (fi *FileIdentifier) IsPlainText(path string) bool {
	extensions := []string{
		".txt",
		".html",
		".js",
		"jsx",
		".ts",
		".tsx",
		".css",
		".sass",
		".scss",
		".go",
		".py",
		".rb",
		".java",
		".c",
		".h",
		".cpp",
		".hpp",
		".json",
		".yml",
		".yaml",
		".toml",
		".md",
	}
	extension := filepath.Ext(path)
	for _, v := range extensions {
		if strings.ToLower(extension) == v {
			return true
		}
	}
	return false
}

func (fi *FileIdentifier) IsImage(path string) bool {
	extensions := []string{
		".xpm",
		".png",
		".jpg",
		".jpeg",
		".jp2",
		".gif",
		".webp",
		".tiff",
		".bmp",
		".ico",
		".heif",
		".xcf",
		".svg",
	}
	extension := filepath.Ext(path)
	for _, v := range extensions {
		if strings.ToLower(extension) == v {
			return true
		}
	}
	return false
}

func (fi *FileIdentifier) IsNonAlphaChannelImage(path string) bool {
	extensions := []string{
		".jpg",
		".jpeg",
		".gif",
		".tiff",
		".bmp",
	}
	extension := filepath.Ext(path)
	for _, v := range extensions {
		if strings.ToLower(extension) == v {
			return true
		}
	}
	return false
}

func (fi *FileIdentifier) IsVideo(path string) bool {
	extensions := []string{
		".ogv",
		".mpeg",
		".mov",
		".mqv",
		".mp4",
		".webm",
		".3gp",
		".3g2",
		".avi",
		".flv",
		".mkv",
		".asf",
		".m4v",
	}
	extension := filepath.Ext(path)
	for _, v := range extensions {
		if strings.ToLower(extension) == v {
			return true
		}
	}
	return false
}
