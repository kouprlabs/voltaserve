package core

type PipelineOptions struct {
	FileID     string `json:"fileId"`
	SnapshotID string `json:"snapshotId"`
	Bucket     string `json:"bucket"`
	Key        string `json:"key"`
}

type PipelineResponse struct {
	Preview   *S3Object
	Text      *S3Object
	OCR       *S3Object
	Thumbnail *Thumbnail
}

type S3Object struct {
	Bucket string      `json:"bucket,omitempty"`
	Key    string      `json:"key,omitempty"`
	Size   int64       `json:"size"`
	Image  *ImageProps `json:"image,omitempty"`
}

type ImageProps struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Thumbnail struct {
	Base64 string `json:"base64"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
