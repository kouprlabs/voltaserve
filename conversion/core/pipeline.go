package core

type Pipeline interface {
	Run(PipelineOptions) error
}
