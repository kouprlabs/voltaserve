export function getSizeWithAspectRatio(
  srcWidth: number,
  srcHeight: number,
  maxWidth: number,
  maxHeight: number,
) {
  const ratio = Math.min(maxWidth / srcWidth, maxHeight / srcHeight)
  return { width: srcWidth * ratio, height: srcHeight * ratio }
}
