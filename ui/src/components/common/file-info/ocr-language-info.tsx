import { useEffect, useState } from 'react'
import { Badge, Stat, StatLabel, StatNumber, Text } from '@chakra-ui/react'
import { Spinner, variables } from '@koupr/ui'
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

  if (!file.ocr?.language) {
    return null
  }

  return (
    <Stat>
      <StatLabel>OCR Language</StatLabel>
      {language ? (
        <StatNumber fontSize={variables.bodyFontSize}>
          <Badge>{language}</Badge>
        </StatNumber>
      ) : (
        <Spinner />
      )}
    </Stat>
  )
}

export default OcrLanguageInfo
