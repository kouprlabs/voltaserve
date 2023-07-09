import { useEffect, useState } from 'react'
import { Badge, Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { File } from '@/client/api/file'
import OcrLanguageAPI from '@/client/api/ocr-language'

type OcrLanguageInfoProps = {
  file: File
}

const OcrLanguageInfo = ({ file }: OcrLanguageInfoProps) => {
  const [language, setLanguage] = useState<string>()

  useEffect(() => {
    async function fetch(language: string) {
      const result = await OcrLanguageAPI.list({ query: language })
      setLanguage(result.data[0].name)
    }
    if (file.ocr?.language) {
      fetch(file.ocr?.language)
    }
  }, [file])

  if (!language) {
    return null
  }

  return (
    <Stat>
      <StatLabel>OCR Language</StatLabel>
      <StatNumber fontSize={variables.bodyFontSize}>
        <Badge>{language}</Badge>
      </StatNumber>
    </Stat>
  )
}

export default OcrLanguageInfo
