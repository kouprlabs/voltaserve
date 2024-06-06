package core

type ToolRunOptions struct {
	Bin    string   `json:"bin"`
	Args   []string `json:"args"`
	Stdout bool     `json:"stdout"`
}

type PipelineRunOptions struct {
	PipelineID *string  `json:"pipelineId"`
	SnapshotID string   `json:"snapshotId"`
	Bucket     string   `json:"bucket"`
	Key        string   `json:"key"`
	Values     []string `json:"values,omitempty"`
}

type SnapshotUpdateOptions struct {
	Options   PipelineRunOptions `json:"options,omitempty"`
	Original  *S3Object          `json:"original,omitempty"`
	Preview   *S3Object          `json:"preview,omitempty"`
	Text      *S3Object          `json:"text,omitempty"`
	OCR       *S3Object          `json:"ocr,omitempty"`
	Entities  *S3Object          `json:"entities,omitempty"`
	Mosaic    *S3Object          `json:"mosaic,omitempty"`
	Watermark *S3Object          `json:"watermark,omitempty"`
	Thumbnail *ImageBase64       `json:"thumbnail,omitempty"`
	Status    string             `json:"status,omitempty"`
}

type S3Object struct {
	Bucket string      `json:"bucket"`
	Key    string      `json:"key"`
	Size   *int64      `json:"size,omitempty"`
	Image  *ImageProps `json:"image,omitempty"`
}

type ImageProps struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type ImageBase64 struct {
	Base64 string `json:"base64"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Pipeline interface {
	Run(PipelineRunOptions) error
}

type Builder interface {
	Build(PipelineRunOptions) error
}
