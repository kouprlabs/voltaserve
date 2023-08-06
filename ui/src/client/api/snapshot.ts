export type Snapshot = {
  version: number
  original: Download
  preview?: Download
  text?: Download
  thumbnail?: Thumbnail
}

export type Download = {
  extension: string
  size: number
  image?: ImageProps
}

export type ImageProps = {
  width: number
  height: number
}

export type Thumbnail = {
  base64: string
  width: number
  height: number
}
