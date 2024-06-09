package processor

import (
	"voltaserve/infra"
)

type ZIPProcessor struct {
	cmd *infra.Command
}

func NewZIPProcessor() *ZIPProcessor {
	return &ZIPProcessor{
		cmd: infra.NewCommand(),
	}
}

func (p *ZIPProcessor) Extract(inputPath string, outputDir string) error {
	if err := p.cmd.Exec("unzip", inputPath, "-d", outputDir); err != nil {
		return err
	}
	return nil
}
