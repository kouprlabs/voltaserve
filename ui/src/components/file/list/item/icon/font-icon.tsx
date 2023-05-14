import { useMemo } from 'react'
import {
  FaFilePdf,
  FaFileExcel,
  FaFileWord,
  FaFilePowerpoint,
  FaFileAudio,
  FaFileVideo,
  FaFileCode,
  FaFileCsv,
  FaFileArchive,
  FaFileAlt,
  FaFile,
  FaFont,
  FaFileContract,
  FaFileImage,
} from 'react-icons/fa'
import { File } from '@/api/file'
import * as fileExtension from '@/helpers/file-extension'

type FontIconProps = {
  file: File
  scale: number
}

const SIZE = 72

const FontIcon = ({ file, scale }: FontIconProps) => {
  const size = useMemo(() => `${SIZE * scale}px`, [scale])
  const isImage = useMemo(
    () =>
      file.original?.extension &&
      fileExtension.isImage(file.original.extension),
    [file.original]
  )
  const isPdf = useMemo(
    () =>
      file.original?.extension && fileExtension.isPdf(file.original.extension),
    [file.original]
  )
  const isText = useMemo(
    () =>
      file.original?.extension && fileExtension.isText(file.original.extension),
    [file.original]
  )
  const isRichText = useMemo(
    () =>
      file.original?.extension &&
      fileExtension.isRichText(file.original.extension),
    [file.original]
  )
  const isWord = useMemo(
    () =>
      file.original?.extension && fileExtension.isWord(file.original.extension),
    [file.original]
  )
  const isPowerPoint = useMemo(
    () =>
      file.original?.extension &&
      fileExtension.isPowerPoint(file.original.extension),
    [file.original]
  )
  const isExcel = useMemo(
    () =>
      file.original?.extension &&
      fileExtension.isExcel(file.original.extension),
    [file.original]
  )
  const isDocument = useMemo(
    () =>
      file.original?.extension &&
      fileExtension.isDocument(file.original.extension),
    [file.original]
  )
  const isSpreadsheet = useMemo(
    () =>
      file.original?.extension &&
      fileExtension.isSpreadsheet(file.original.extension),
    [file.original]
  )
  const isSlides = useMemo(
    () =>
      file.original?.extension &&
      fileExtension.isSlides(file.original.extension),
    [file.original]
  )
  const isVideo = useMemo(
    () =>
      file.original?.extension &&
      fileExtension.isVideo(file.original.extension),
    [file.original]
  )
  const isAudio = useMemo(
    () =>
      file.original?.extension &&
      fileExtension.isAudio(file.original.extension),
    [file.original]
  )
  const isArchive = useMemo(
    () =>
      file.original?.extension &&
      fileExtension.isArchive(file.original.extension),
    [file.original]
  )
  const isFont = useMemo(
    () =>
      file.original?.extension && fileExtension.isFont(file.original.extension),
    [file.original]
  )
  const isCode = useMemo(
    () =>
      file.original?.extension && fileExtension.isCode(file.original.extension),
    [file.original]
  )
  const isCSV = useMemo(
    () =>
      file.original?.extension && fileExtension.isCSV(file.original.extension),
    [file.original]
  )

  if (isPdf) {
    return <FaFilePdf fontSize={size} />
  }
  if (isImage) {
    return <FaFileImage fontSize={size} />
  }
  if (isText) {
    return <FaFileAlt fontSize={size} />
  }
  if (isWord) {
    return <FaFileWord fontSize={size} />
  }
  if (isRichText) {
    return <FaFileContract fontSize={size} />
  }
  if (isDocument) {
    return <FaFileWord fontSize={size} />
  }
  if (isExcel) {
    return <FaFileExcel fontSize={size} />
  }
  if (isSpreadsheet) {
    return <FaFileExcel fontSize={size} />
  }
  if (isPowerPoint) {
    return <FaFilePowerpoint fontSize={size} />
  }
  if (isSlides) {
    return <FaFilePowerpoint fontSize={size} />
  }
  if (isArchive) {
    return <FaFileArchive fontSize={size} />
  }
  if (isFont) {
    return <FaFont fontSize={size} />
  }
  if (isAudio) {
    return <FaFileAudio fontSize={size} />
  }
  if (isVideo) {
    return <FaFileVideo fontSize={size} />
  }
  if (isCode) {
    return <FaFileCode fontSize={size} />
  }
  if (isCSV) {
    return <FaFileCsv fontSize={size} />
  }
  return <FaFile fontSize={size} />
}

export default FontIcon
