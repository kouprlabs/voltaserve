import { File } from '@/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

type AudioPlayerProps = {
  file: File
}

const AudioPlayer = ({ file }: AudioPlayerProps) => {
  if (!file.original) {
    return null
  }
  return (
    <audio controls>
      <source
        src={`/proxy/api/v1/files/${file.id}/original${
          file.original.extension
        }?${new URLSearchParams({
          access_token: getAccessTokenOrRedirect(),
        })}`}
      />
    </audio>
  )
}

export default AudioPlayer
