export type Snapshot = {
  version: number
  original: Download
  preview?: Download
  ocr?: Download
  text?: Download
  thumbnail?: string
}

export type Download = {
  extension: string
  size: number
  image: ImageProps | undefined
}

export type ImageProps = {
  width: number
  height: number
}
