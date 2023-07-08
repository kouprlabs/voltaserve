package core

type PipelineOptions struct {
	FileID                string `json:"fileId"`
	SnapshotID            string `json:"snapshotId"`
	Bucket                string `json:"bucket"`
	Key                   string `json:"key"`
	IsAutomaticOCREnabled bool   `json:"isAutomaticOcrEnabled"`
}

type PipelineResponse struct {
	Options   PipelineOptions `json:"options,omitempty"`
	Original  *S3Object       `json:"original,omitempty"`
	Preview   *S3Object       `json:"preview,omitempty"`
	Text      *S3Object       `json:"text,omitempty"`
	OCR       *S3Object       `json:"ocr,omitempty"`
	Thumbnail *ImageBase64    `json:"thumbnail,omitempty"`
}

type S3Object struct {
	Bucket   string      `json:"bucket"`
	Key      string      `json:"key"`
	Size     int64       `json:"size"`
	Image    *ImageProps `json:"image,omitempty"`
	Language *string     `json:"language,omitempty"`
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
	Run(PipelineOptions) error
}

type Builder interface {
	Build(PipelineOptions) error
}
