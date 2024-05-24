import { File } from '@/client/api/file'
import { getSizeWithAspectRatio } from '@/helpers/aspect-ratio'

const MAX_WIDTH = 130
const MAX_HEIGHT = 130

export function getThumbnailWidth(file: File, scale: number): string {
  if (file.snapshot?.thumbnail) {
    const { width } = getSizeWithAspectRatio(
      file.snapshot?.thumbnail.width,
      file.snapshot?.thumbnail.height,
      MAX_WIDTH,
      MAX_HEIGHT,
    )
    return `${width * scale}px`
  } else {
    return `${MAX_WIDTH * scale}px`
  }
}

export function getThumbnailHeight(file: File, scale: number): string {
  if (file.snapshot?.thumbnail) {
    const { height } = getSizeWithAspectRatio(
      file.snapshot?.thumbnail.width,
      file.snapshot?.thumbnail.height,
      MAX_WIDTH,
      MAX_HEIGHT,
    )
    return `${height * scale}px`
  } else {
    return `${MAX_HEIGHT * scale}px`
  }
}
