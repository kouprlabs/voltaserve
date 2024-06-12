import { useMemo, useState } from 'react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import { leViewerPermission } from '@/client/api/permission'
import { getAccessTokenOrRedirect } from '@/infra/token'
import SectionSpinner from '@/lib/components/section-spinner'
import variables from '@/lib/variables'

export type ViewerImageProps = {
  file: File
}

const ViewerImage = ({ file }: ViewerImageProps) => {
  const [isLoading, setIsLoading] = useState(true)
  const isWatermark = useMemo(
    () => location.pathname.endsWith('/watermark'),
    [location],
  )
  const url = useMemo(() => {
    if (
      file.snapshot?.watermark?.extension &&
      (isWatermark || leViewerPermission(file.permission))
    ) {
      return `/proxy/api/v2/watermarks/${file.id}/watermark${
        file.snapshot?.watermark.extension
      }?${new URLSearchParams({
        access_token: getAccessTokenOrRedirect(),
      })}`
    } else if (file.snapshot?.preview && file.snapshot?.preview.extension) {
      return `/proxy/api/v2/files/${file.id}/preview${
        file.snapshot?.preview.extension
      }?${new URLSearchParams({
        access_token: getAccessTokenOrRedirect(),
      })}`
    }
  }, [file])

  if (!file.snapshot?.preview) {
    return null
  }

  return (
    <div className={cx('flex', 'flex-col', 'w-full', 'h-full', 'gap-1.5')}>
      <div
        className={cx(
          'relative',
          'flex',
          'items-center',
          'justify-center',
          'grow',
          'w-full',
          'h-full',
          'overflow-scroll',
        )}
      >
        {isLoading && <SectionSpinner />}
        <img
          src={url}
          style={{
            objectFit: 'contain',
            width: isLoading ? 0 : 'auto',
            height: isLoading ? 0 : '90%',
            visibility: isLoading ? 'hidden' : 'visible',
            borderRadius: variables.borderRadius,
          }}
          onLoad={() => setIsLoading(false)}
          alt={file.name}
        />
      </div>
    </div>
  )
}

export default ViewerImage
