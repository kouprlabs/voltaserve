export type Tile = {
  width: number
  height: number
  lastColWidth: number
  lastRowHeight: number
}

export type ZoomLevel = {
  index: number
  width: number
  height: number
  rows: number
  cols: number
  scaleDownPercentage: number
  tile: Tile
}

export type Metadata = {
  width: number
  height: number
  extension: string
  zoomLevels: ZoomLevel[]
}
