package model

type WatermarkInfo struct {
	IsAvailable bool               `json:"isAvailable"`
	Metadata    *WatermarkMetadata `json:"metadata,omitempty"`
}

type WatermarkMetadata struct {
	IsOutdated bool `json:"isOutdated"`
}
