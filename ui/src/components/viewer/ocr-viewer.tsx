import { useMemo } from 'react'
import { File } from '@/client/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

type OCRViewerProps = {
  file: File
}

const OCRViewer = ({ file }: OCRViewerProps) => {
  const download = useMemo(() => file.ocr, [file])
  const url = useMemo(() => {
    if (!download || !download.extension) {
      return ''
    }
    return `/proxy/api/v1/files/${file.id}/ocr${
      download.extension
    }?${new URLSearchParams({
      access_token: getAccessTokenOrRedirect(),
    })}`
  }, [file, download])

  if (!download) {
    return null
  }

  return <iframe width="100%" height="100%" src={url} title={file.name} />
}

export default OCRViewer
