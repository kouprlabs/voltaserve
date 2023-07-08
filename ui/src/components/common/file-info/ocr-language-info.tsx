import { Badge, Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { File } from '@/client/api/file'

type OcrLanguageInfoProps = {
  file: File
}

const OcrLanguageInfo = ({ file }: OcrLanguageInfoProps) => {
  if (!file.ocr?.language) {
    return null
  }
  return (
    <Stat>
      <StatLabel>OCR Language</StatLabel>
      <StatNumber fontSize={variables.bodyFontSize}>
        <Badge>{file.ocr.language}</Badge>
      </StatNumber>
    </Stat>
  )
}

export default OcrLanguageInfo
