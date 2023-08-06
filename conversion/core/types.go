package core

type PipelineRunOptions struct {
	FileID     string `json:"fileId"`
	SnapshotID string `json:"snapshotId"`
	Bucket     string `json:"bucket"`
	Key        string `json:"key"`
}

type SnapshotUpdateOptions struct {
	Options   PipelineRunOptions `json:"options,omitempty"`
	Original  *S3Object          `json:"original,omitempty"`
	Preview   *S3Object          `json:"preview,omitempty"`
	Text      *S3Object          `json:"text,omitempty"`
	Thumbnail *ImageBase64       `json:"thumbnail,omitempty"`
	Status    string             `json:"status,omitempty"`
}

type S3Object struct {
	Bucket string      `json:"bucket"`
	Key    string      `json:"key"`
	Size   int64       `json:"size"`
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
