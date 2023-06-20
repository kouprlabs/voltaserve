package infra

import (
	"os"
	"path/filepath"
	"strings"
	"voltaserve/helper"
)

type OfficeProcessor struct {
	cmd *Command
}

func NewOfficeProcessor() *OfficeProcessor {
	return &OfficeProcessor{
		cmd: NewCommand(),
	}
}

func (p *OfficeProcessor) PDF(inputPath string) (string, error) {
	outputDir := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}
	if err := p.cmd.Exec("soffice", "--headless", "--convert-to", "pdf", inputPath, "--outdir", outputDir); err != nil {
		return "", err
	}
	filename := filepath.Base(inputPath)
	outputPath := filepath.FromSlash(outputDir + "/" + strings.TrimSuffix(filename, filepath.Ext(filename)) + ".pdf")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return "", err
	}
	return outputPath, nil
}
