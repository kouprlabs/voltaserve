package processor

import (
	"voltaserve/infra"
)

type GLTFProcessor struct {
	cmd *infra.Command
}

func NewGLTFProcessor() *GLTFProcessor {
	return &GLTFProcessor{
		cmd: infra.NewCommand(),
	}
}

func (p *GLTFProcessor) ToGLB(inputPath string, outputPath string) error {
	if err := p.cmd.Exec("gltf-pipeline", "-i", inputPath, "-o", outputPath); err != nil {
		return err
	}
	return nil
}
