export function isPdf(extension: string) {
  return extension === '.pdf'
}

export function isImage(extension: string) {
  return (
    [
      '.xpm',
      '.png',
      '.jpg',
      '.jpeg',
      '.jp2',
      '.gif',
      '.webp',
      '.tiff',
      '.bmp',
      '.ico',
      '.heif',
      '.xcf',
      '.svg',
    ].findIndex((e) => e === extension) !== -1
  )
}

export function isText(extension: string) {
  return extension === '.txt'
}

export function isRichText(extension: string) {
  return ['.rtf'].findIndex((e) => e === extension) !== -1
}

export function isWord(extension: string) {
  return ['.docx', '.doc'].findIndex((e) => e === extension) !== -1
}

export function isPowerPoint(extension: string) {
  return ['.pptx', '.ppt'].findIndex((e) => e === extension) !== -1
}

export function isExcel(extension: string) {
  return ['.xlsx', '.xls'].findIndex((e) => e === extension) !== -1
}

export function isDocument(extension: string) {
  return (
    ['.odt', '.ott', '.gdoc', '.pages'].findIndex((e) => e === extension) !== -1
  )
}

export function isSpreadsheet(extension: string) {
  return ['.ods', '.ots', '.gsheet'].findIndex((e) => e === extension) !== -1
}

export function isSlides(extension: string) {
  return (
    ['.odp', '.otp', '.key', '.gslides'].findIndex((e) => e === extension) !==
    -1
  )
}

export function isVideo(extension: string) {
  return (
    [
      '.ogv',
      '.mpeg',
      '.mov',
      '.mqv',
      '.mp4',
      '.webm',
      '.3gp',
      '.3g2',
      '.avi',
      '.flv',
      '.mkv',
      '.asf',
      '.m4v',
    ].findIndex((e) => e === extension) !== -1
  )
}

export function isAudio(extension: string) {
  return (
    [
      '.oga',
      '.ogg',
      '.mp3',
      '.flac',
      '.midi',
      '.ape',
      '.mpc',
      '.amr',
      '.wav',
      '.aiff',
      '.au',
      '.aac',
      'voc',
      '.m4a',
      '.qcp',
    ].findIndex((e) => e === extension) !== -1
  )
}

export function isArchive(extension: string) {
  return (
    ['.zip', '.tar', '.7z', '.bz2', '.gz', '.rar'].findIndex(
      (e) => e === extension
    ) !== -1
  )
}

export function isFont(extension: string) {
  return ['.ttf', '.woff'].findIndex((e) => e === extension) !== -1
}

export function isCode(extension: string) {
  return (
    [
      '.html',
      '.js',
      'jsx',
      '.ts',
      '.tsx',
      '.css',
      '.sass',
      '.scss',
      '.go',
      '.py',
      '.rb',
      '.java',
      '.c',
      '.h',
      '.cpp',
      '.hpp',
      '.json',
      '.yml',
      '.yaml',
      '.toml',
      '.md',
    ].findIndex((e) => e === extension) !== -1
  )
}

export function isCSV(extension: string) {
  return ['.csv'].findIndex((e) => e === extension) !== -1
}
