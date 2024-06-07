package model

import "voltaserve/client"

type Pipeline interface {
	Run(client.PipelineRunOptions) error
}

type Builder interface {
	Build(client.PipelineRunOptions) error
}
