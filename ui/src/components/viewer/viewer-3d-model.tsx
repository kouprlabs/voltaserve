import { useMemo } from 'react'
import '@google/model-viewer'
import { File } from '@/client/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

export type Viewer3DModelProps = {
  file: File
}

const Viewer3DModel = ({ file }: Viewer3DModelProps) => {
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
      {/* @ts-ignore */}
      <model-viewer
        src={url}
        shadow-intensity="1"
        camera-controls
        touch-action="pan-y"
        style={{ width: '100%', height: '100%' }}
      >
        {/* @ts-ignore */}
      </model-viewer>
    </>
  )
}

export default Viewer3DModel
