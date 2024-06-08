package processor

import (
	"voltaserve/config"
	"voltaserve/infra"
)

type GLTFProcessor struct {
	cmd    *infra.Command
	config config.Config
}

func NewGLTFProcessor() *GLTFProcessor {
	return &GLTFProcessor{
		cmd:    infra.NewCommand(),
		config: config.GetConfig(),
	}
}

func (p *GLTFProcessor) ToGLB(inputPath string, outputPath string) error {
	if err := p.cmd.Exec("gltf-pipeline", "-i", inputPath, "-o", outputPath); err != nil {
		return err
	}
	return nil
}
