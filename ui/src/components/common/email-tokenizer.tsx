import { useMemo } from 'react'
import { Tag } from '@chakra-ui/react'
import cx from 'classnames'
import parseEmailList from '@/helpers/parse-email-list'

export type EmailTokenizerProps = {
  value: string
}

const EmailTokenizer = ({ value }: EmailTokenizerProps) => {
  const emails = useMemo(() => parseEmailList(value), [value])
  return (
    <>
      {emails.length > 0 ? (
        <div className={cx('flex', 'flex-wrap', 'gap-0.5')}>
          {emails.map((email, index) => (
            <Tag
              key={index}
              size="md"
              variant="solid"
              className={cx('rounded-full')}
            >
              {email}
            </Tag>
          ))}
        </div>
      ) : null}
    </>
  )
}

export default EmailTokenizer
