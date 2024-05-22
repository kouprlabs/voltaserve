import { useMemo } from 'react'
import { File } from '@/client/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

export type ViewerVideoProps = {
  file: File
}

const ViewerVideo = ({ file }: ViewerVideoProps) => {
  const download = useMemo(() => file.snapshot?.original, [file])
  const url = useMemo(() => {
    if (!download || !download.extension) {
      return ''
    }
    return `/proxy/api/v2/files/${file.id}/original${
      download.extension
    }?${new URLSearchParams({
      access_token: getAccessTokenOrRedirect(),
    })}`
  }, [file, download])

  if (!download) {
    return null
  }

  return (
    <video controls autoPlay style={{ maxWidth: '100%', maxHeight: '100%' }}>
      <source src={url} />
    </video>
  )
}

export default ViewerVideo
