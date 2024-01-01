import { useEffect, useRef } from 'react'
import { Text } from '@chakra-ui/react'

interface TruncateTextProps {
  text: string
  maxCharacters: number
}

const TruncatedText = ({ text, maxCharacters }: TruncateTextProps) => {
  const elementRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (elementRef.current && text.length > maxCharacters) {
      elementRef.current.textContent = text.slice(0, maxCharacters).trim() + '…'
    }
  }, [text, maxCharacters])

  return (
    <Text as="span" ref={elementRef} whiteSpace="nowrap" overflow="hidden">
      {text}
    </Text>
  )
}

export default TruncatedText
