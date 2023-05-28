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
	outputDirectory := filepath.FromSlash(os.TempDir() + "/" + helper.NewId())
	if err := os.MkdirAll(outputDirectory, 0755); err != nil {
		return "", err
	}
	if err := p.cmd.Exec("soffice", "--headless", "--convert-to", "pdf", inputPath, "--outdir", outputDirectory); err != nil {
		return "", err
	}
	filename := filepath.Base(inputPath)
	outputPath := filepath.FromSlash(outputDirectory + "/" + strings.TrimSuffix(filename, filepath.Ext(filename)) + ".pdf")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return "", err
	}
	return outputPath, nil
}
