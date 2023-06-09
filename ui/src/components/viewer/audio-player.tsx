import { useMemo } from 'react'
import { File } from '@/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

type AudioPlayerProps = {
  file: File
}

const AudioPlayer = ({ file }: AudioPlayerProps) => {
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
    <audio controls>
      <source src={url} />
    </audio>
  )
}

export default AudioPlayer
