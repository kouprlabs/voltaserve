import { useColorMode } from '@chakra-ui/react'
import { File } from '@/client/api/file'
import * as fe from '@/helpers/file-extension'
import DarkArchiveSvg from './assets/dark-icon-archive.svg'
import DarkAudioSvg from './assets/dark-icon-audio.svg'
import DarkCodeSvg from './assets/dark-icon-code.svg'
import DarkCsvSvg from './assets/dark-icon-csv.svg'
import DarkFileSvg from './assets/dark-icon-file.svg'
import DarkPdfSvg from './assets/dark-icon-pdf.svg'
import DarkPowerPointSvg from './assets/dark-icon-power-point.svg'
import DarkRichTextSvg from './assets/dark-icon-rich-text.svg'
import DarkSpreadsheetSvg from './assets/dark-icon-spreadsheet.svg'
import DarkTextSvg from './assets/dark-icon-text.svg'
import DarkVideoSvg from './assets/dark-icon-video.svg'
import DarkWordSvg from './assets/dark-icon-word.svg'
import ArchiveSvg from './assets/icon-archive.svg'
import AudioSvg from './assets/icon-audio.svg'
import CodeSvg from './assets/icon-code.svg'
import CsvSvg from './assets/icon-csv.svg'
import FileSvg from './assets/icon-file.svg'
import PdfSvg from './assets/icon-pdf.svg'
import PowerPointSvg from './assets/icon-power-point.svg'
import RichTextSvg from './assets/icon-rich-text.svg'
import SpreadsheetSvg from './assets/icon-spreadsheet.svg'
import TextSvg from './assets/icon-text.svg'
import VideoSvg from './assets/icon-video.svg'
import WordSvg from './assets/icon-word.svg'

export type IconFontProps = {
  file: File
  scale: number
}

const MIN_WIDTH = 45
const MIN_HEIGHT = 59.78
const BASE_WIDTH = 67
const BASE_HEIGHT = 89

const IconDiverse = ({ file, scale }: IconFontProps) => {
  const { colorMode } = useColorMode()
  const width = Math.max(MIN_WIDTH, BASE_WIDTH * scale)
  const height = Math.max(MIN_HEIGHT, BASE_HEIGHT * scale)

  const { original } = file
  let image
  if (fe.isImage(original?.extension)) {
    image = colorMode === 'dark' ? DarkFileSvg : FileSvg
  } else if (fe.isPDF(original?.extension)) {
    image = colorMode === 'dark' ? DarkPdfSvg : PdfSvg
  } else if (fe.isText(original?.extension)) {
    image = colorMode === 'dark' ? DarkTextSvg : TextSvg
  } else if (fe.isRichText(original?.extension)) {
    image = colorMode === 'dark' ? DarkRichTextSvg : RichTextSvg
  } else if (fe.isWord(original?.extension)) {
    image = colorMode === 'dark' ? DarkWordSvg : WordSvg
  } else if (fe.isPowerPoint(original?.extension)) {
    image = colorMode === 'dark' ? DarkPowerPointSvg : PowerPointSvg
  } else if (fe.isExcel(original?.extension)) {
    image = colorMode === 'dark' ? DarkSpreadsheetSvg : SpreadsheetSvg
  } else if (fe.isDocument(original?.extension)) {
    image = colorMode === 'dark' ? DarkWordSvg : WordSvg
  } else if (fe.isSpreadsheet(original?.extension)) {
    image = colorMode === 'dark' ? DarkSpreadsheetSvg : SpreadsheetSvg
  } else if (fe.isSlides(original?.extension)) {
    image = colorMode === 'dark' ? DarkPowerPointSvg : PowerPointSvg
  } else if (fe.isVideo(original?.extension)) {
    image = colorMode === 'dark' ? DarkVideoSvg : VideoSvg
  } else if (fe.isAudio(original?.extension)) {
    image = colorMode === 'dark' ? DarkAudioSvg : AudioSvg
  } else if (fe.isArchive(original?.extension)) {
    image = colorMode === 'dark' ? DarkArchiveSvg : ArchiveSvg
  } else if (fe.isFont(original?.extension)) {
    image = colorMode === 'dark' ? DarkFileSvg : FileSvg
  } else if (fe.isCode(original?.extension)) {
    image = colorMode === 'dark' ? DarkCodeSvg : CodeSvg
  } else if (fe.isCSV(original?.extension)) {
    image = colorMode === 'dark' ? DarkCsvSvg : CsvSvg
  } else {
    image = colorMode === 'dark' ? DarkFileSvg : FileSvg
  }

  return (
    <img src={image} style={{ width: `${width}px`, height: `${height}px` }} />
  )
}

export default IconDiverse
