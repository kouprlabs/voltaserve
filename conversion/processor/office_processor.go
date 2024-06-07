package processor

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"voltaserve/config"
	"voltaserve/helper"
	"voltaserve/infra"
)

type OfficeProcessor struct {
	cmd    *infra.Command
	config config.Config
}

func NewOfficeProcessor() *OfficeProcessor {
	return &OfficeProcessor{
		cmd:    infra.NewCommand(),
		config: config.GetConfig(),
	}
}

func (p *OfficeProcessor) PDF(inputPath string) (*string, error) {
	outputDir := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
	if err := infra.NewCommand().Exec("soffice", "--headless", "--convert-to", "pdf", "--outdir", outputDir, inputPath); err != nil {
		return nil, err
	}
	if _, err := os.Stat(inputPath); err != nil {
		return nil, err
	}
	base := filepath.Base(inputPath)
	return helper.ToPtr(filepath.FromSlash(outputDir + "/" + strings.TrimSuffix(base, path.Ext(base)) + ".pdf")), nil
}
