package core

type Pipeline interface {
	Run(PipelineOptions) (PipelineResponse, error)
}
