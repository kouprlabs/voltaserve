import { useMemo } from 'react'
import { File } from '@/api/file'

type PdfViewerProps = {
  file: File
}

const PdfViewer = ({ file }: PdfViewerProps) => {
  const url = useMemo(
    () => `/proxy/api/v1/files/${file.id}/preview${file!.preview!.extension}`,
    [file]
  )
  if (!file.preview) {
    return null
  }
  return <iframe width="100%" height="100%" src={url} title={file.name} />
}

export default PdfViewer
