package helper

import (
	"fmt"
	"voltaserve/core"

	"github.com/fatih/color"
)

func SprintPipelineOptions(opts *core.PipelineOptions) string {
	return fmt.Sprintf("(FileID=%s SnapshotID=%s S3Bucket=%s S3Key=%s)", opts.FileID, opts.SnapshotID, opts.Bucket, opts.Key)
}

func PrintPipelineOptions(opts *core.PipelineOptions) {
	color.Set(color.FgHiBlack)
	fmt.Print(SprintPipelineOptions(opts))
	color.Unset()
}

func PrintlnPipelineOptions(opts *core.PipelineOptions) {
	PrintPipelineOptions(opts)
	fmt.Println()
}
