import { Button } from '@chakra-ui/react'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import { swrConfig } from '@/client/options'
import { IconOpenInNew } from '@/lib/components/icons'
import { useAppSelector } from '@/store/hook'

const WatermarkOverviewArtifacts = () => {
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const { data: file } = FileAPI.useGet(id, swrConfig())

  if (!file) {
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
      <Button
        as="a"
        type="button"
        leftIcon={<IconOpenInNew />}
        target="_blank"
        href={`/file/${file.id}/watermark`}
      >
        Open Watermark-Protected File
      </Button>
    </div>
  )
}

export default WatermarkOverviewArtifacts
