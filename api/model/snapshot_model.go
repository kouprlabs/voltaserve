package model

type SnapshotModel interface {
	GetID() string
	GetVersion() int64
	GetOriginal() *S3Object
	GetPreview() *S3Object
	GetText() *S3Object
	GetOCR() *S3Object
	GetThumbnail() *Thumbnail
	HasOriginal() bool
	HasPreview() bool
	HasText() bool
	HasOCR() bool
	HasThumbnail() bool
	GetCreateTime() string
	GetUpdateTime() *string
	SetOriginal(*S3Object)
	SetPreview(*S3Object)
	SetText(*S3Object)
	SetOCR(*S3Object)
	SetThumbnail(*Thumbnail)
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
