export function isPDF(ext?: string | null) {
  if (!ext) {
    return false
  }
  return ext === '.pdf'
}

export function isImage(ext?: string | null) {
  if (!ext) {
    return false
  }
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
    ].findIndex((e) => e === ext) !== -1
  )
}

export function isText(ext?: string | null) {
  if (!ext) {
    return false
  }
  return ext === '.txt'
}

export function isRichText(ext?: string | null) {
  if (!ext) {
    return false
  }
  return ['.rtf'].findIndex((e) => e === ext) !== -1
}

export function isWord(ext?: string | null) {
  if (!ext) {
    return false
  }
  return ['.docx', '.doc'].findIndex((e) => e === ext) !== -1
}

export function isPowerPoint(ext?: string | null) {
  if (!ext) {
    return false
  }
  return ['.pptx', '.ppt'].findIndex((e) => e === ext) !== -1
}

export function isExcel(ext?: string | null) {
  if (!ext) {
    return false
  }
  return ['.xlsx', '.xls'].findIndex((e) => e === ext) !== -1
}

export function isMicrosoftOffice(ext?: string | null) {
  return isWord(ext) || isPowerPoint(ext) || isExcel(ext)
}

export function isOpenOffice(ext?: string | null) {
  return isDocument(ext) || isSpreadsheet(ext) || isSlides(ext)
}

export function isDocument(ext?: string | null) {
  if (!ext) {
    return false
  }
  return ['.odt', '.ott', '.gdoc', '.pages'].findIndex((e) => e === ext) !== -1
}

export function isSpreadsheet(ext?: string | null) {
  if (!ext) {
    return false
  }
  return ['.ods', '.ots', '.gsheet'].findIndex((e) => e === ext) !== -1
}

export function isSlides(ext?: string | null) {
  if (!ext) {
    return false
  }
  return ['.odp', '.otp', '.key', '.gslides'].findIndex((e) => e === ext) !== -1
}

export function isVideo(ext?: string | null) {
  if (!ext) {
    return false
  }
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
    ].findIndex((e) => e === ext) !== -1
  )
}

export function isAudio(ext?: string | null) {
  if (!ext) {
    return false
  }
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
    ].findIndex((e) => e === ext) !== -1
  )
}

export function isArchive(ext?: string | null) {
  if (!ext) {
    return false
  }
  return (
    ['.zip', '.tar', '.7z', '.bz2', '.gz', '.rar'].findIndex(
      (e) => e === ext,
    ) !== -1
  )
}

export function isFont(ext?: string | null) {
  if (!ext) {
    return false
  }
  return ['.ttf', '.woff'].findIndex((e) => e === ext) !== -1
}

export function isCode(ext?: string | null) {
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
    ].findIndex((e) => e === ext) !== -1
  )
}

export function isCSV(ext?: string | null) {
  if (!ext) {
    return false
  }
  return ['.csv'].findIndex((e) => e === ext) !== -1
}

export function is3DModel(ext?: string | null) {
  if (!ext) {
    return false
  }
  return ext === '.glb'
}
