import { useMemo } from 'react'
import { File } from '@/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

type VideoPlayerProps = {
  file: File
}

const VideoPlayer = ({ file }: VideoPlayerProps) => {
  const url = useMemo(() => {
    if (file.original?.extension) {
      const searchParams = new URLSearchParams({
        access_token: getAccessTokenOrRedirect(),
      })
      return `/proxy/api/v1/files/${file.id}/original${file.original.extension}?${searchParams}`
    } else {
      return ''
    }
  }, [file])
  return (
    <video controls autoPlay style={{ maxWidth: '100%', maxHeight: '100%' }}>
      <source src={url} />
    </video>
  )
}

export default VideoPlayer
