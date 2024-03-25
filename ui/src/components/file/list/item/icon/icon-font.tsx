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
import { File } from '@/client/api/file'
import * as fe from '@/helpers/file-extension'

type IconFontProps = {
  file: File
  scale: number
}

const SIZE = 72

const IconFont = ({ file, scale }: IconFontProps) => {
  const size = `${SIZE * scale}px`
  const { original } = file

  if (fe.isImage(original?.extension)) {
    return <FaFileImage fontSize={size} />
  } else if (fe.isPDF(original?.extension)) {
    return <FaFilePdf fontSize={size} />
  } else if (fe.isText(original?.extension)) {
    return <FaFileAlt fontSize={size} />
  } else if (fe.isRichText(original?.extension)) {
    return <FaFileContract fontSize={size} />
  } else if (fe.isWord(original?.extension)) {
    return <FaFileWord fontSize={size} />
  } else if (fe.isPowerPoint(original?.extension)) {
    return <FaFilePowerpoint fontSize={size} />
  } else if (fe.isExcel(original?.extension)) {
    return <FaFileExcel fontSize={size} />
  } else if (fe.isDocument(original?.extension)) {
    return <FaFileWord fontSize={size} />
  } else if (fe.isSpreadsheet(original?.extension)) {
    return <FaFileExcel fontSize={size} />
  } else if (fe.isSlides(original?.extension)) {
    return <FaFilePowerpoint fontSize={size} />
  } else if (fe.isVideo(original?.extension)) {
    return <FaFileVideo fontSize={size} />
  } else if (fe.isAudio(original?.extension)) {
    return <FaFileAudio fontSize={size} />
  } else if (fe.isArchive(original?.extension)) {
    return <FaFileArchive fontSize={size} />
  } else if (fe.isFont(original?.extension)) {
    return <FaFont fontSize={size} />
  } else if (fe.isCode(original?.extension)) {
    return <FaFileCode fontSize={size} />
  } else if (fe.isCSV(original?.extension)) {
    return <FaFileCsv fontSize={size} />
  } else {
    return <FaFile fontSize={size} />
  }
}

export default IconFont
