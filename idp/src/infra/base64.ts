export function base64ToBuffer(value: string): Buffer {
  let withoutPrefix: string
  if (value.includes(',')) {
    withoutPrefix = value.split(',')[1]
  } else {
    withoutPrefix = value
  }
  try {
    return Buffer.from(withoutPrefix, 'base64')
  } catch (err) {
    throw new Error(err as string)
  }
}

export function base64ToMIME(value: string): string {
  if (!value.startsWith('data:image/')) {
    return ''
  }
  const colonIndex = value.indexOf(':')
  const semicolonIndex = value.indexOf(';')
  if (colonIndex === -1 || semicolonIndex === -1) {
    return ''
  }
  return value.substring(colonIndex + 1, semicolonIndex)
}

export function base64ToExtension(value: string): string {
  const mime = base64ToMIME(value)
  switch (mime) {
    case 'image/jpeg':
      return '.jpg'
    case 'image/png':
      return '.png'
    default:
      return ''
  }
}