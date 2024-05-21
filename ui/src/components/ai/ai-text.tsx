import { Button } from '@chakra-ui/react'
import cx from 'classnames'
import AIAPI from '@/client/api/analysis'
import { getAccessTokenOrRedirect } from '@/infra/token'
import { IconOpenInNew } from '@/lib'

export type AITextProps = {
  id: string
}

const AIText = ({ id }: AITextProps) => {
  const { data: summary } = AIAPI.useGetSummary(id)

  if (!summary) {
    return null
  }

  return (
    <div
      className={cx(
        'flex',
        'flex-col',
        'items-center',
        'justify-center',
        'gap-1',
      )}
    >
      {summary.hasText ? (
        <Button
          as="a"
          type="button"
          leftIcon={<IconOpenInNew />}
          href={`/proxy/api/v2/files/${id}/text.txt?${new URLSearchParams({
            access_token: getAccessTokenOrRedirect(),
          })}`}
          target="_blank"
        >
          Open Text
        </Button>
      ) : null}
      {summary.hasOcr ? (
        <Button
          as="a"
          type="button"
          leftIcon={<IconOpenInNew />}
          href={`/proxy/api/v2/files/${id}/ocr.pdf?${new URLSearchParams({
            access_token: getAccessTokenOrRedirect(),
          })}`}
          target="_blank"
        >
          Open Searchable PDF
        </Button>
      ) : null}
    </div>
  )
}

export default AIText
