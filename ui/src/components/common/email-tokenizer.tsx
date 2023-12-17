import { useMemo } from 'react'
import { Tag, Wrap, WrapItem } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import * as Yup from 'yup'

export function parseEmailList(value: string): string[] {
  return [...new Set(value.split(',').map((e: string) => e.trim()))].filter(
    (e) => {
      if (e.length === 0) {
        return false
      }
      try {
        Yup.string()
          .email()
          .matches(
            /.+(\.[A-Za-z]{2,})$/,
            'Email must end with a valid top-level domain',
          )
          .validateSync(e)
        return true
      } catch {
        return false
      }
    },
  )
}

type EmailTokenizerProps = {
  value: string
}

const EmailTokenizer = ({ value }: EmailTokenizerProps) => {
  const emails = useMemo(() => parseEmailList(value), [value])
  return (
    <>
      {emails.length > 0 ? (
        <Wrap spacing={variables.spacingXs}>
          {emails.map((email, index) => (
            <WrapItem key={index}>
              <Tag size="md" borderRadius="full" variant="solid">
                {email}
              </Tag>
            </WrapItem>
          ))}
        </Wrap>
      ) : null}
    </>
  )
}

export default EmailTokenizer
