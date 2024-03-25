import { useMemo, useState } from 'react'
import { SectionSpinner } from '@koupr/ui'
import classNames from 'classnames'
import { File } from '@/client/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

export type ViewerImageProps = {
  file: File
}

const ViewerImage = ({ file }: ViewerImageProps) => {
  const [isLoading, setIsLoading] = useState(true)
  const download = useMemo(() => file.preview ?? file.original, [file])
  const path = useMemo(() => (file.preview ? 'preview' : 'original'), [file])
  const url = useMemo(() => {
    if (!download?.extension) {
      return ''
    }
    return `/proxy/api/v1/files/${file.id}/${path}${
      download.extension
    }?${new URLSearchParams({
      access_token: getAccessTokenOrRedirect(),
    })}`
  }, [file, download, path])

  if (!download) {
    return null
  }

  return (
    <div
      className={classNames('flex', 'flex-col', 'w-full', 'h-full', 'gap-1.5')}
    >
      <div
        className={classNames(
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
