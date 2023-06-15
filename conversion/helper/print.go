package helper

import (
	"fmt"
	"voltaserve/core"
)

func SprintPipelineOptions(opts *core.PipelineOptions) string {
	return fmt.Sprintf("(FileID=%s SnapshotID=%s S3Bucket=%s S3Key=%s)", opts.FileID, opts.SnapshotID, opts.Bucket, opts.Key)
}
