import { useMemo } from 'react'
import { File } from '@/client/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

export type ViewerPDFProps = {
  file: File
}

const ViewerPDF = ({ file }: ViewerPDFProps) => {
  const download = useMemo(() => file.preview || file.original, [file])
  const urlPath = useMemo(() => (file.preview ? 'preview' : 'original'), [file])
  const url = useMemo(() => {
    if (!download || !download.extension) {
      return ''
    }
    return `/proxy/api/v1/files/${file.id}/${urlPath}${
      download.extension
    }?${new URLSearchParams({
      access_token: getAccessTokenOrRedirect(),
    })}`
  }, [file, download, urlPath])

  if (!download) {
    return null
  }

  return <iframe width="100%" height="100%" src={url} title={file.name} />
}

export default ViewerPDF
