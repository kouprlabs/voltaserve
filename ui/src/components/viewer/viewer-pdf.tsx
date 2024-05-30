import { useMemo } from 'react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import { leViewerPermission } from '@/client/api/permission'
import { getAccessTokenOrRedirect } from '@/infra/token'

export type ViewerPDFProps = {
  file: File
}

const ViewerPDF = ({ file }: ViewerPDFProps) => {
  const isWatermark = useMemo(
    () => location.pathname.endsWith('/watermark'),
    [location],
  )
  const url = useMemo(() => {
    if (
      file.snapshot?.watermark?.extension &&
      (isWatermark || leViewerPermission(file.permission))
    ) {
      return `/proxy/api/v2/watermarks/${file.id}/watermark${
        file.snapshot?.watermark.extension
      }?${new URLSearchParams({
        access_token: getAccessTokenOrRedirect(),
      })}`
    } else if (file.snapshot?.preview && file.snapshot?.preview.extension) {
      return `/proxy/api/v2/files/${file.id}/preview${
        file.snapshot?.preview.extension
      }?${new URLSearchParams({
        access_token: getAccessTokenOrRedirect(),
      })}`
    }
  }, [file])

  if (!file.snapshot?.preview) {
    return null
  }

  return (
    <iframe className={cx('w-full', 'h-full')} src={url} title={file.name} />
  )
}

export default ViewerPDF
