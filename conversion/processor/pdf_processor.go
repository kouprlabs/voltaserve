// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package processor

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/infra"
)

type PDFProcessor struct {
	cmd       *infra.Command
	imageProc *ImageProcessor
	config    *config.Config
}

func NewPDFProcessor() *PDFProcessor {
	return &PDFProcessor{
		cmd:       infra.NewCommand(),
		imageProc: NewImageProcessor(),
		config:    config.GetConfig(),
	}
}

func (p *PDFProcessor) TextFromPDF(inputPath string) (*string, error) {
	tmpPath := filepath.Join(os.TempDir(), helper.NewID()+".txt")

	if err := infra.NewCommand().Exec("pdftotext", inputPath, tmpPath); err != nil {
		return nil, err
	}

	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			infra.GetLogger().Error(err)
		}
	}(tmpPath)

	b, err := os.ReadFile(tmpPath) //nolint:gosec // Known path
	if err != nil {
		return nil, err
	}

	return helper.ToPtr(strings.TrimSpace(string(b))), nil
}

func (p *PDFProcessor) Thumbnail(inputPath string, width int, height int, outputPath string) error {
	var widthStr string
	if width == 0 {
		widthStr = ""
	} else {
		widthStr = strconv.FormatInt(int64(width), 10)
	}
	var heightStr string
	if height == 0 {
		heightStr = ""
	} else {
		heightStr = strconv.FormatInt(int64(height), 10)
	}
	if err := infra.NewCommand().Exec("convert", "-thumbnail", widthStr+"x"+heightStr, "-background", "white", "-alpha", "remove", "-flatten", fmt.Sprintf("%s[0]", inputPath), outputPath); err != nil {
		return err
	}
	return nil
}

func (p *PDFProcessor) SplitPages(inputPath string, outputDir string) error {
	if err := infra.NewCommand().Exec("qpdf", "--split-pages", inputPath, filepath.FromSlash(outputDir+"/%d.pdf")); err != nil {
		return err
	}
	/* Rename files by removing leading zeros */
	if err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".pdf") {
			base := filepath.Base(path)
			ext := filepath.Ext(base)
			re := regexp.MustCompile(`(\D*)(\d+)(.*)`)
			matches := re.FindStringSubmatch(strings.TrimSuffix(base, ext))
			if len(matches) == 4 {
				number, err := strconv.Atoi(matches[2])
				if err != nil {
					return err
				}
				newName := fmt.Sprintf("%s%d%s%s", matches[1], number, matches[3], ext)
				newPath := filepath.Join(filepath.Dir(path), newName)
				if err = os.Rename(path, newPath); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (p *PDFProcessor) SplitThumbnails(inputPath string, width int, height int, extension string, outputDir string) error {
	var widthStr string
	if width == 0 {
		widthStr = ""
	} else {
		widthStr = strconv.FormatInt(int64(width), 10)
	}
	var heightStr string
	if height == 0 {
		heightStr = ""
	} else {
		heightStr = strconv.FormatInt(int64(height), 10)
	}
	dimensions := widthStr + "x" + heightStr
	if err := infra.NewCommand().Exec("convert", "-thumbnail", dimensions, "-gravity", "center", "-extent", dimensions, inputPath, fmt.Sprintf("%s%s", filepath.FromSlash(outputDir+"/%d"), extension)); err != nil {
		return err
	}

	// Rename all files in outputDir to add +1 to their filenames.
	// For example: 0.jpg, 1.jpg, 2.jpg become 1.jpg, 2.jpg and 3.jpg respectively.
	// We start the renaming from the last file (largest number) to the smallest (0.jpg).
	// We consider all content of outputDir to be a flat hierarchy of files, no recursive processing is needed.

	files, err := os.ReadDir(outputDir)
	if err != nil {
		return err
	}

	// Create a list of files with their integer values.
	var structs []struct {
		Name  string
		Index int
	}
	for _, file := range files {
		baseName := file.Name()
		index, err := strconv.Atoi(strings.TrimSuffix(baseName, filepath.Ext(baseName)))
		if err != nil {
			continue
		}
		structs = append(structs, struct {
			Name  string
			Index int
		}{
			Name:  baseName,
			Index: index,
		})
	}

	// Sort the file names by their index in descending order.
	sort.Slice(structs, func(i, j int) bool {
		return structs[i].Index > structs[j].Index
	})

	// Rename files to add +1 to their index.
	for _, file := range structs {
		oldPath := filepath.Join(outputDir, file.Name)
		newPath := filepath.Join(outputDir, fmt.Sprintf("%d%s", file.Index+1, filepath.Ext(file.Name)))
		if err := os.Rename(oldPath, newPath); err != nil {
			return err
		}
	}

	return nil
}

func (p *PDFProcessor) CountPages(inputPath string) (*int, error) {
	output, err := infra.NewCommand().ReadOutput("qpdf", "--show-npages", inputPath)
	if err != nil {
		return nil, err
	}
	count, err := strconv.Atoi(strings.TrimSpace(*output))
	if err != nil {
		return nil, err
	}
	return &count, nil
}
