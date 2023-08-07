package processor

import (
	"go.uber.org/zap"
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
	logger *zap.SugaredLogger
}

func NewOfficeProcessor() *OfficeProcessor {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &OfficeProcessor{
		cmd:    infra.NewCommand(),
		config: config.GetConfig(),
		logger: logger,
	}
}

func (p *OfficeProcessor) PDF(inputPath string) (string, error) {
	outputDir := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
	if err := infra.NewCommand().Exec("soffice", "--headless", "--convert-to", "pdf", "--outdir", outputDir, inputPath); err != nil {
		return "", err
	}
	if _, err := os.Stat(inputPath); err != nil {
		return "", err
	}
	base := filepath.Base(inputPath)
	return filepath.FromSlash(outputDir + "/" + strings.TrimSuffix(base, path.Ext(base)) + ".pdf"), nil
}
