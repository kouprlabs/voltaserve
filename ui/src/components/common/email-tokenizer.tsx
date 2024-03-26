import { useMemo } from 'react'
import { Tag, Wrap, WrapItem } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import classNames from 'classnames'
import parseEmailList from '@/helpers/parse-email-list'

export type EmailTokenizerProps = {
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
              <Tag
                size="md"
                variant="solid"
                className={classNames('rounded-full')}
              >
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
