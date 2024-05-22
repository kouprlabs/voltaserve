import { useMemo } from 'react'
import { File } from '@/client/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

export type ViewerAudioProps = {
  file: File
}

const ViewerAudio = ({ file }: ViewerAudioProps) => {
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
    <audio controls>
      <source src={url} />
    </audio>
  )
}

export default ViewerAudio
