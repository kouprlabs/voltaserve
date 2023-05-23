import { File, FileType } from '@/api/file'
import { SortDirection, SortType } from '@/models/sort'
import {
  isArchive,
  isAudio,
  isCode,
  isCSV,
  isDocument,
  isExcel,
  isFont,
  isImage,
  isPDF,
  isPowerPoint,
  isRichText,
  isSlides,
  isSpreadsheet,
  isText,
  isVideo,
  isWord,
} from './file-extension'

export function sortByName(data: File[], lt: number, gt: number): File[] {
  data.sort((a, b) => {
    if (a.name.toLowerCase() < b.name.toLowerCase()) {
      return lt
    }
    if (a.name.toLowerCase() > b.name.toLowerCase()) {
      return gt
    }
    return 0
  })
  return data
}

export function sortByKind(data: File[], direction: SortDirection): File[] {
  const folders = data.filter((e) => e.type === FileType.Folder)
  let files = data.filter((e) => e.type === FileType.File)
  const images = data.filter(
    (e) => e.original?.extension && isImage(e.original.extension)
  )
  const pdfs = data.filter(
    (e) => e.original?.extension && isPDF(e.original.extension)
  )
  const documents = data.filter(
    (e) => e.original?.extension && isDocument(e.original.extension)
  )
  const words = data.filter(
    (e) => e.original?.extension && isWord(e.original.extension)
  )
  const spreadsheets = data.filter(
    (e) => e.original?.extension && isSpreadsheet(e.original.extension)
  )
  const excels = data.filter(
    (e) => e.original?.extension && isExcel(e.original.extension)
  )
  const slides = data.filter(
    (e) => e.original?.extension && isSlides(e.original.extension)
  )
  const powerpoints = data.filter(
    (e) => e.original?.extension && isPowerPoint(e.original.extension)
  )
  const videos = data.filter(
    (e) => e.original?.extension && isVideo(e.original.extension)
  )
  const audios = data.filter(
    (e) => e.original?.extension && isAudio(e.original.extension)
  )
  const archives = data.filter(
    (e) => e.original?.extension && isArchive(e.original.extension)
  )
  const texts = data.filter(
    (e) => e.original?.extension && isText(e.original.extension)
  )
  const richTexts = data.filter(
    (e) => e.original?.extension && isRichText(e.original.extension)
  )
  const csvs = data.filter(
    (e) => e.original?.extension && isCSV(e.original.extension)
  )
  const codes = data.filter(
    (e) => e.original?.extension && isCode(e.original.extension)
  )
  const fonts = data.filter(
    (e) => e.original?.extension && isFont(e.original.extension)
  )
  const others = data.filter(
    (e) =>
      e.original?.extension &&
      !isImage(e.original.extension) &&
      !isPDF(e.original.extension) &&
      !isDocument(e.original.extension) &&
      !isWord(e.original.extension) &&
      !isSpreadsheet(e.original.extension) &&
      !isExcel(e.original.extension) &&
      !isSlides(e.original.extension) &&
      !isPowerPoint(e.original.extension) &&
      !isVideo(e.original.extension) &&
      !isAudio(e.original.extension) &&
      !isArchive(e.original.extension) &&
      !isText(e.original.extension) &&
      !isRichText(e.original.extension) &&
      !isCSV(e.original.extension) &&
      !isCode(e.original.extension) &&
      !isFont(e.original.extension)
  )
  files = [
    ...images,
    ...pdfs,
    ...documents,
    ...words,
    ...spreadsheets,
    ...excels,
    ...slides,
    ...powerpoints,
    ...videos,
    ...audios,
    ...archives,
    ...texts,
    ...richTexts,
    ...csvs,
    ...codes,
    ...fonts,
    ...others,
  ]
  if (direction === SortDirection.Ascending) {
    return [...folders, ...files]
  } else {
    return [...files, ...folders]
  }
}

export function sortBySize(data: File[], lt: number, gt: number): File[] {
  data.sort((a, b) => {
    const sizeA = a.original?.size || 0
    const sizeB = b.original?.size || 0
    if (sizeA < sizeB) {
      return lt
    }
    if (sizeA > sizeB) {
      return gt
    }
    return 0
  })
  return data
}

export function sortByDateCreated(
  data: File[],
  lt: number,
  gt: number
): File[] {
  data.sort((a, b) => {
    const dateA = new Date(a.createTime)
    const dateB = new Date(b.createTime)
    if (dateA < dateB) {
      return lt
    }
    if (dateA > dateB) {
      return gt
    }
    return 0
  })
  return data
}

export function sortByDateModified(
  data: File[],
  lt: number,
  gt: number
): File[] {
  data.sort((a, b) => {
    const dateA = new Date(a.updateTime || a.createTime)
    const dateB = new Date(b.updateTime || b.createTime)
    if (dateA < dateB) {
      return lt
    }
    if (dateA > dateB) {
      return gt
    }
    return 0
  })
  return data
}

export function sort(
  data: File[],
  type: SortType,
  direction: SortDirection
): File[] {
  const lt = direction === SortDirection.Ascending ? -1 : 1
  const gt = direction === SortDirection.Ascending ? 1 : -1
  if (type === SortType.ByName) {
    return sortByName(data, lt, gt)
  } else if (type === SortType.ByKind) {
    return sortByKind(data, direction)
  } else if (type === SortType.BySize) {
    return sortBySize(data, lt, gt)
  } else if (type === SortType.ByDateCreated) {
    return sortByDateCreated(data, lt, gt)
  } else if (type === SortType.ByDateModified) {
    return sortByDateModified(data, lt, gt)
  } else {
    throw new Error('Invalid sort type')
  }
}
