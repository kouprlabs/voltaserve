import { useMemo, useState } from 'react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'
import { SectionSpinner } from '@/lib'

export type ViewerImageProps = {
  file: File
}

const ViewerImage = ({ file }: ViewerImageProps) => {
  const [isLoading, setIsLoading] = useState(true)
  const url = useMemo(() => {
    if (!file.preview?.extension) {
      return ''
    }
    return `/proxy/api/v1/files/${file.id}/preview${
      file.preview.extension
    }?${new URLSearchParams({
      access_token: getAccessTokenOrRedirect(),
    })}`
  }, [file])

  if (!file.preview) {
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
            height: isLoading ? 0 : '100%',
            visibility: isLoading ? 'hidden' : 'visible',
          }}
          onLoad={() => setIsLoading(false)}
          alt={file.name}
        />
      </div>
    </div>
  )
}

export default ViewerImage
