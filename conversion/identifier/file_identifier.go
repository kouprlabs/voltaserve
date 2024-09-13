// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package identifier

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
)

type FileIdentifier struct{}

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
		".pages",
		".numbers",
		".key",
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
		".csv",
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
		".tif",
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
		".tif",
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

func (fi *FileIdentifier) IsKRA(path string) (bool, error) {
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return false, fmt.Errorf("open zip: %w", err)
	}

	defer zipFile.Close()

	mimetypeFile, err := zipFile.Open("mimetype")
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("open mimetype: %w", err)
	}
	defer mimetypeFile.Close()

	mimetype, err := io.ReadAll(mimetypeFile)
	if err != nil {
		return false, fmt.Errorf("read mimetype: %w", err)
	}

	return string(mimetype) == "application/x-krita", nil
}
