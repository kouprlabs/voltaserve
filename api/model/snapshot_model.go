package model

const (
	SnapshotStatusNew        = "new"
	SnapshotStatusProcessing = "processing"
	SnapshotStatusReady      = "ready"
	SnapshotStatusError      = "error"
)

type Snapshot interface {
	GetID() string
	GetVersion() int64
	GetOriginal() *S3Object
	GetPreview() *S3Object
	GetText() *S3Object
	GetOCR() *S3Object
	GetThumbnail() *Thumbnail
	GetLanguage() *string
	HasOriginal() bool
	HasPreview() bool
	HasText() bool
	HasOCR() bool
	HasThumbnail() bool
	HasLanguage() bool
	GetStatus() string
	GetCreateTime() string
	GetUpdateTime() *string
	SetID(string)
	SetVersion(int64)
	SetOriginal(*S3Object)
	SetPreview(*S3Object)
	SetText(*S3Object)
	SetOCR(*S3Object)
	SetThumbnail(*Thumbnail)
	SetLanguage(*string)
	SetStatus(string)
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

type Thumbnail struct {
	Base64 string `json:"base64"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
