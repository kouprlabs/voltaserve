import { useMemo } from 'react'
import { File } from '@/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

type VideoPlayerProps = {
  file: File
}

const VideoPlayer = ({ file }: VideoPlayerProps) => {
  const download = useMemo(() => file.original, [file])
  const url = useMemo(() => {
    if (!download || !download.extension) {
      return ''
    }
    return `/proxy/api/v1/files/${file.id}/original${
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

export default VideoPlayer
