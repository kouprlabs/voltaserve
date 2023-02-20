import { useMemo } from 'react'
import { Box } from '@chakra-ui/react'
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
} from 'react-icons/fa'
import { File } from '@/api/file'
import * as fileExtension from '@/helpers/file-extension'
import { ItemSize } from '..'
import FileListItemImageIcon from './image-icon'
import FileListItemSharedSign from './shared-sign'

type FileListItemFileIconProps = {
  file: File
  size: ItemSize
}

const FileListItemFileIcon = ({ file, size }: FileListItemFileIconProps) => {
  const fontSize = useMemo(() => {
    if (size === 'normal') {
      return '72px'
    }
    if (size === 'large') {
      return '150px'
    }
  }, [size])
  const isPdf = useMemo(
    () =>
      file.original?.extension && fileExtension.isPdf(file.original.extension),
    [file.original]
  )
  const isImage = useMemo(
    () =>
      file.original?.extension &&
      fileExtension.isImage(file.original.extension),
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

  const renderIcon = () => {
    if (isPdf) {
      return <FaFilePdf fontSize={fontSize} />
    }
    if (isText) {
      return <FaFileAlt fontSize={fontSize} />
    }
    if (isWord) {
      return <FaFileWord fontSize={fontSize} />
    }
    if (isRichText) {
      return <FaFileContract fontSize={fontSize} />
    }
    if (isDocument) {
      return <FaFileWord fontSize={fontSize} />
    }
    if (isExcel) {
      return <FaFileExcel fontSize={fontSize} />
    }
    if (isSpreadsheet) {
      return <FaFileExcel fontSize={fontSize} />
    }
    if (isPowerPoint) {
      return <FaFilePowerpoint fontSize={fontSize} />
    }
    if (isSlides) {
      return <FaFilePowerpoint fontSize={fontSize} />
    }
    if (isArchive) {
      return <FaFileArchive fontSize={fontSize} />
    }
    if (isFont) {
      return <FaFont fontSize={fontSize} />
    }
    if (isAudio) {
      return <FaFileAudio fontSize={fontSize} />
    }
    if (isVideo) {
      return <FaFileVideo fontSize={fontSize} />
    }
    if (isCode) {
      return <FaFileCode fontSize={fontSize} />
    }
    if (isCSV) {
      return <FaFileCsv fontSize={fontSize} />
    }
    return <FaFile fontSize={fontSize} />
  }

  if (isImage) {
    return <FileListItemImageIcon file={file} size={size} />
  } else {
    return (
      <Box position="relative">
        {renderIcon()}
        {file.isShared && <FileListItemSharedSign bottom="-5px" right="0px" />}
      </Box>
    )
  }
}

export default FileListItemFileIcon
