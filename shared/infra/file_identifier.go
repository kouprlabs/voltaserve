// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package infra

import (
	"archive/zip"
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"github.com/kouprlabs/voltaserve/shared/logger"
)

type FileIdentifier struct{}

func NewFileIdentifier() *FileIdentifier {
	return &FileIdentifier{}
}

func (fi *FileIdentifier) IsPDF(path string) bool {
	return strings.ToLower(filepath.Ext(path)) == ".pdf"
}

func (fi *FileIdentifier) IsOffice(path string) bool {
	for _, v := range []string{
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
	} {
		if strings.ToLower(filepath.Ext(path)) == v {
			return true
		}
	}
	return false
}

func (fi *FileIdentifier) IsPlainText(path string) bool {
	for _, v := range []string{
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
	} {
		if strings.ToLower(filepath.Ext(path)) == v {
			return true
		}
	}
	return false
}

func (fi *FileIdentifier) IsDocument(path string) bool {
	return fi.IsPDF(path) || fi.IsOffice(path) || fi.IsPlainText(path)
}

func (fi *FileIdentifier) IsImage(path string) bool {
	return fi.IsJPEG(path) ||
		fi.IsPNG(path) ||
		fi.IsTIFF(path) ||
		fi.IsGIF(path) ||
		fi.IsWebP(path) ||
		fi.IsXPM(path) ||
		fi.IsBMP(path) ||
		fi.IsICO(path) ||
		fi.IsSVG(path) ||
		fi.IsHEIF(path) ||
		fi.IsXCF(path)
}

func (fi *FileIdentifier) IsJPEG(path string) bool {
	for _, v := range []string{
		".jpg",
		".jpeg",
		".jpe",
		".jfif",
		".jif",
		".jp2",
	} {
		if strings.ToLower(filepath.Ext(path)) == v {
			return true
		}
	}
	return false
}

func (fi *FileIdentifier) IsPNG(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".png"
}

func (fi *FileIdentifier) IsTIFF(path string) bool {
	for _, v := range []string{
		".tiff",
		".tif",
	} {
		if strings.ToLower(filepath.Ext(path)) == v {
			return true
		}
	}
	return false
}

func (fi *FileIdentifier) IsGIF(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".gif"
}

func (fi *FileIdentifier) IsWebP(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".webp"
}

func (fi *FileIdentifier) IsXPM(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".xpm"
}

func (fi *FileIdentifier) IsBMP(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".bmp"
}

func (fi *FileIdentifier) IsICO(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".ico"
}

func (fi *FileIdentifier) IsSVG(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".svg"
}

func (fi *FileIdentifier) IsHEIF(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".heif"
}

func (fi *FileIdentifier) IsXCF(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".xcf"
}

func (fi *FileIdentifier) IsNonAlphaChannelImage(path string) bool {
	for _, v := range []string{
		".jpg",
		".jpeg",
		".gif",
		".tiff",
		".tif",
		".bmp",
	} {
		if strings.ToLower(filepath.Ext(path)) == v {
			return true
		}
	}
	return false
}

func (fi *FileIdentifier) IsVideo(path string) bool {
	for _, v := range []string{
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
	} {
		if strings.ToLower(filepath.Ext(path)) == v {
			return true
		}
	}
	return false
}

func (fi *FileIdentifier) IsAudio(path string) bool {
	for _, v := range []string{
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
	} {
		if strings.ToLower(filepath.Ext(path)) == v {
			return true
		}
	}
	return false
}

func (fi *FileIdentifier) IsGLB(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".glb"
}

func (fi *FileIdentifier) IsOBJ(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".obj"
}

func (fi *FileIdentifier) IsFBX(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".fbx"
}

func (fi *FileIdentifier) IsSTL(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".stl"
}

func (fi *FileIdentifier) IsPLY(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".ply"
}

func (fi *FileIdentifier) IsSTEP(path string) bool {
	for _, v := range []string{
		".step",
		".stp",
	} {
		if strings.ToLower(filepath.Ext(path)) == v {
			return true
		}
	}
	return false
}

func (fi *FileIdentifier) Is3DS(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".3ds"
}

func (fi *FileIdentifier) IsBLEND(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".blend"
}

func (fi *FileIdentifier) IsMAX(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".max"
}

func (fi *FileIdentifier) IsC4D(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".c4d"
}

func (fi *FileIdentifier) Is3D(path string) bool {
	return fi.IsGLB(path) ||
		fi.IsOBJ(path) ||
		fi.IsFBX(path) ||
		fi.IsSTL(path) ||
		fi.IsPLY(path) ||
		fi.IsSTEP(path) ||
		fi.Is3DS(path) ||
		fi.IsBLEND(path) ||
		fi.IsMAX(path) ||
		fi.IsC4D(path)
}

func (fi *FileIdentifier) IsZIP(path string) bool {
	for _, v := range []string{
		".zip",
		".zipx",
	} {
		if strings.ToLower(filepath.Ext(path)) == v {
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

// IsGLTF inspects a ZIP archive to see if it contains a glTF 2.0 structure.
func (fi *FileIdentifier) IsGLTF(path string) (bool, error) {
	if !fi.IsZIP(path) {
		return false, nil
	}
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return false, err
	}
	defer func(zipFile *zip.ReadCloser) {
		if err := zipFile.Close(); err != nil {
			logger.GetLogger().Error(err)
		}
	}(zipFile)
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
			defer func(rc io.ReadCloser) {
				if err := rc.Close(); err != nil {
					logger.GetLogger().Error(err)
				}
			}(rc)
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
	return hasGLTF && (!hasBin || (gltfFile != nil)), nil
}
