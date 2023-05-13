import { useMemo, useState } from 'react'
import {
  Box,
  Image,
  Skeleton,
  useColorModeValue,
  useToken,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
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
import ImageIcon from './image-icon'
import SharedSign from './shared-sign'
import { getThumbnailHeight, getThumbnailWidth } from './thumbnail-size'

type FileIconProps = {
  file: File
  scale: number
}

const ICON_FONT_SIZE = 72

const FileIcon = ({ file, scale }: FileIconProps) => {
  const [isLoading, setIsLoading] = useState(true)
  const width = useMemo(() => getThumbnailWidth(file, scale), [scale, file])
  const height = useMemo(() => getThumbnailHeight(file, scale), [scale, file])
  const iconFontSize = useMemo(() => {
    return `${ICON_FONT_SIZE * scale}px`
  }, [scale])
  const borderColor = useColorModeValue('gray.300', 'gray.700')
  const [borderColorDecoded] = useToken('colors', [borderColor])
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
      return <FaFilePdf fontSize={iconFontSize} />
    }
    if (isText) {
      return <FaFileAlt fontSize={iconFontSize} />
    }
    if (isWord) {
      return <FaFileWord fontSize={iconFontSize} />
    }
    if (isRichText) {
      return <FaFileContract fontSize={iconFontSize} />
    }
    if (isDocument) {
      return <FaFileWord fontSize={iconFontSize} />
    }
    if (isExcel) {
      return <FaFileExcel fontSize={iconFontSize} />
    }
    if (isSpreadsheet) {
      return <FaFileExcel fontSize={iconFontSize} />
    }
    if (isPowerPoint) {
      return <FaFilePowerpoint fontSize={iconFontSize} />
    }
    if (isSlides) {
      return <FaFilePowerpoint fontSize={iconFontSize} />
    }
    if (isArchive) {
      return <FaFileArchive fontSize={iconFontSize} />
    }
    if (isFont) {
      return <FaFont fontSize={iconFontSize} />
    }
    if (isAudio) {
      return <FaFileAudio fontSize={iconFontSize} />
    }
    if (isVideo) {
      return <FaFileVideo fontSize={iconFontSize} />
    }
    if (isCode) {
      return <FaFileCode fontSize={iconFontSize} />
    }
    if (isCSV) {
      return <FaFileCsv fontSize={iconFontSize} />
    }
    return <FaFile fontSize={iconFontSize} />
  }

  if (isImage) {
    return <ImageIcon file={file} scale={scale} />
  } else {
    if (file.thumbnail) {
      return (
        <Box position="relative" width={width} height={height}>
          <Image
            src={file.thumbnail?.base64}
            width={isLoading ? 0 : width}
            height={isLoading ? 0 : height}
            style={{
              objectFit: 'cover',
              width: isLoading ? 0 : width,
              height: isLoading ? 0 : height,
              border: '1px solid',
              borderColor: borderColorDecoded,
              borderRadius: variables.borderRadiusSm,
              visibility: isLoading ? 'hidden' : 'visible',
            }}
            alt={file.name}
            onLoad={() => setIsLoading(false)}
          />
          {isLoading && (
            <Skeleton
              width={width}
              height={height}
              borderRadius={variables.borderRadiusSm}
            />
          )}
          {file.isShared && <SharedSign bottom="-5px" right="-5px" />}
        </Box>
      )
    } else {
      return (
        <Box position="relative">
          {renderIcon()}
          {file.isShared && <SharedSign bottom="-5px" right="0px" />}
        </Box>
      )
    }
  }
}

export default FileIcon
