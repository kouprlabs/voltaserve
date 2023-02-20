import { File } from '@/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

type VideoPlayerProps = {
  file: File
}

const VideoPlayer = ({ file }: VideoPlayerProps) => (
  <video controls autoPlay style={{ maxWidth: '100%', maxHeight: '100%' }}>
    <source
      src={`/proxy/api/v1/files/${file.id}/original${
        file.original!.extension
      }?${new URLSearchParams({
        access_token: getAccessTokenOrRedirect(),
      })}`}
    />
  </video>
)

export default VideoPlayer
