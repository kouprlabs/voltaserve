package conversion

type PipelineOptions struct {
	FileID     string
	SnapshotID string
	S3Bucket   string
	S3Key      string
}

type Pipeline interface {
	Run(opts PipelineOptions) error
}
