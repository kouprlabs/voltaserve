import { Button } from '@chakra-ui/react'
import cx from 'classnames'
import InsightsAPI from '@/client/api/insights'
import { getAccessTokenOrRedirect } from '@/infra/token'
import { IconOpenInNew } from '@/lib'
import { useAppSelector } from '@/store/hook'

const InsightsOverviewText = () => {
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const { data: summary } = InsightsAPI.useGetSummary(id)

  if (!id || !summary) {
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
          href={`/proxy/api/v2/insights/${id}/text.txt?${new URLSearchParams({
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
          href={`/proxy/api/v2/insights/${id}/ocr.pdf?${new URLSearchParams({
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

export default InsightsOverviewText