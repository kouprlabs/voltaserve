package model

type MosaicMetadata struct {
	Width      int               `json:"width"`
	Height     int               `json:"height"`
	Extension  string            `json:"extension"`
	ZoomLevels []MosaicZoomLevel `json:"zoomLevels"`
}

type MosaicZoomLevel struct {
	Index               int        `json:"index"`
	Width               int        `json:"width"`
	Height              int        `json:"height"`
	Rows                int        `json:"rows"`
	Cols                int        `json:"cols"`
	ScaleDownPercentage float32    `json:"scaleDownPercentage"`
	Tile                MosaicTile `json:"tile"`
}

type MosaicTile struct {
	Width         int `json:"width"`
	Height        int `json:"height"`
	LastColWidth  int `json:"lastColWidth"`
	LastRowHeight int `json:"lastRowHeight"`
}
