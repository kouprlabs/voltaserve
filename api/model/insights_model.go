package model

type InsightsInfo struct {
	IsAvailable bool              `json:"isAvailable"`
	Metadata    *InsightsMetadata `json:"metadata,omitempty"`
}

type InsightsMetadata struct {
	IsOutdated bool `json:"isOutdated"`
}

type InsightsEntity struct {
	Text      string `json:"text"`
	Label     string `json:"label"`
	Frequency int    `json:"frequency"`
}
