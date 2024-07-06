// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package infra

import (
	"archive/zip"
	"encoding/json"
	"path/filepath"
	"strings"
	"voltaserve/config"
)

type FileIdentifier struct {
	config *config.Config
}

func NewFileIdentifier() *FileIdentifier {
	return &FileIdentifier{
		config: config.GetConfig(),
	}
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

func (fi *FileIdentifier) IsAudio(path string) bool {
	extensions := []string{
		".oga",
		".ogg",
		".mp3",
		".flac",
		".midi",
		".ape",
		".mpc",
		".amr",
		".wav",
		".aiff",
		".au",
		".aac",
		"voc",
		".m4a",
		".qcp",
	}
	extension := filepath.Ext(path)
	for _, v := range extensions {
		if strings.ToLower(extension) == v {
			return true
		}
	}
	return false
}

func (fi *FileIdentifier) IsGLB(path string) bool {
	extensions := []string{
		".glb",
	}
	extension := filepath.Ext(path)
	for _, v := range extensions {
		if strings.ToLower(extension) == v {
			return true
		}
	}
	return false
}

func (fi *FileIdentifier) IsZIP(path string) bool {
	extensions := []string{
		".zip",
		".zipx",
	}
	extension := filepath.Ext(path)
	for _, v := range extensions {
		if strings.ToLower(extension) == v {
			return true
		}
	}
	return false
}

type GLTF struct {
	Buffers []struct {
		URI string `json:"uri"`
	} `json:"buffers"`
}

/* Inspects a ZIP archive to see if it contains a glTF 2.0 structure. */
func (fi *FileIdentifier) IsGLTF(path string) (bool, error) {
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return false, err
	}
	defer zipFile.Close()
	var hasGLTF, hasBin bool
	var gltfFile *zip.File
	for _, file := range zipFile.File {
		if strings.HasSuffix(file.Name, ".gltf") {
			hasGLTF = true
			gltfFile = file
		}
		if strings.HasSuffix(file.Name, ".bin") {
			hasBin = true
		}
	}
	if hasGLTF {
		if gltfFile != nil {
			rc, err := gltfFile.Open()
			if err != nil {
				return false, err
			}
			defer rc.Close()
			var gltf GLTF
			if err := json.NewDecoder(rc).Decode(&gltf); err != nil {
				return false, err
			}
			for _, buffer := range gltf.Buffers {
				if strings.HasSuffix(buffer.URI, ".bin") {
					hasBin = true
					break
				}
			}
		}
	}
	return hasGLTF && (!hasBin || (hasBin && gltfFile != nil)), nil
}

func (fi *FileIdentifier) GetProcessingLimitMB(path string) int {
	var res int
	if fi.IsAudio(path) {
		res = fi.config.Limits.GetFileProcessingMB(config.FileTypeAudio)
	} else if fi.IsImage(path) {
		res = fi.config.Limits.GetFileProcessingMB(config.FileTypeImage)
	} else if fi.IsOffice(path) {
		res = fi.config.Limits.GetFileProcessingMB(config.FileTypeOffice)
	} else if fi.IsPDF(path) {
		res = fi.config.Limits.GetFileProcessingMB(config.FileTypePDF)
	} else if fi.IsPlainText(path) {
		res = fi.config.Limits.GetFileProcessingMB(config.FileTypePlainText)
	} else if fi.IsVideo(path) {
		res = fi.config.Limits.GetFileProcessingMB(config.FileTypeVideo)
	} else if fi.IsGLB(path) {
		res = fi.config.Limits.GetFileProcessingMB(config.FileTypeGLB)
	} else if fi.IsZIP(path) {
		res = fi.config.Limits.GetFileProcessingMB(config.FileTypeZIP)
	} else if ok, err := fi.IsGLTF(path); ok && err != nil {
		res = fi.config.Limits.GetFileProcessingMB(config.FileTypeGLTF)
	} else {
		res = fi.config.Limits.GetFileProcessingMB(config.FileTypeEverythingElse)
	}
	return res
}
