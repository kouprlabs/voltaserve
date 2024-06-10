import { useMemo } from 'react'
import '@google/model-viewer'
import { File } from '@/client/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

export type ViewerModelProps = {
  file: File
}

const ViewerModel = ({ file }: ViewerModelProps) => {
  const download = useMemo(() => file.snapshot?.preview, [file])
  const url = useMemo(() => {
    if (!download || !download.extension) {
      return ''
    }
    return `/proxy/api/v2/files/${file.id}/preview${
      download.extension
    }?${new URLSearchParams({
      access_token: getAccessTokenOrRedirect(),
    })}`
  }, [file, download])

  if (!download) {
    return null
  }

  return (
    <>
      {/* @ts-expect-error: ignored */}
      <model-viewer
        src={url}
        shadow-intensity="1"
        camera-controls
        touch-action="pan-y"
        style={{ width: '100%', height: '100%' }}
      >
        {/* @ts-expect-error: ignored */}
      </model-viewer>
    </>
  )
}

export default ViewerModel
