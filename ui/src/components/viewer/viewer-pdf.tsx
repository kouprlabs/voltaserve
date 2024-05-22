import { useMemo } from 'react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

export type ViewerPDFProps = {
  file: File
}

const ViewerPDF = ({ file }: ViewerPDFProps) => {
  const url = useMemo(() => {
    if (!file.snapshot?.preview || !file.snapshot?.preview.extension) {
      return ''
    }
    return `/proxy/api/v2/files/${file.id}/preview${
      file.snapshot?.preview.extension
    }?${new URLSearchParams({
      access_token: getAccessTokenOrRedirect(),
    })}`
  }, [file])

  if (!file.snapshot?.preview) {
    return null
  }

  return (
    <iframe className={cx('w-full', 'h-full')} src={url} title={file.name} />
  )
}

export default ViewerPDF
