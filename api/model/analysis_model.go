package model

type AnalysisEntity struct {
	Text      string `json:"text"`
	Label     string `json:"label"`
	Frequency int    `json:"frequency"`
}
