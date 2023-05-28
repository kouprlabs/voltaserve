package conversion

type PipelineOptions struct {
	FileId     string
	SnapshotId string
	S3Bucket   string
	S3Key      string
}

type Pipeline interface {
	Run(opts PipelineOptions) error
}
