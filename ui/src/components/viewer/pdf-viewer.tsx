import { useMemo } from 'react'
import { File } from '@/api/file'

type PdfViewerProps = {
  file: File
}

const PdfViewer = ({ file }: PdfViewerProps) => {
  const download = useMemo(() => file.preview || file.original, [file])
  const urlPath = useMemo(() => (file.preview ? 'preview' : 'original'), [file])
  const url = useMemo(() => {
    if (!download || !download.extension) {
      return ''
    }
    return `/proxy/api/v1/files/${file.id}/${urlPath}${download.extension}`
  }, [file, download, urlPath])

  if (!download) {
    return null
  }

  return <iframe width="100%" height="100%" src={url} title={file.name} />
}

export default PdfViewer
