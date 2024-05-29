import { File } from '@/client/api/file'
import { Status } from '@/client/api/snapshot'
import IconBadgeError from './icon-badge-error'
import IconBadgeInsights from './icon-badge-insights'
import IconBadgeMosaic from './icon-badge-mosaic'
import IconBadgeNew from './icon-badge-new'
import IconBadgeProcessing from './icon-badge-processing'
import IconBadgeShared from './icon-badge-shared'

export type IconBadgeProps = {
  file: File
  isLoading?: boolean
}

const IconBadge = ({ file, isLoading }: IconBadgeProps) => {
  return (
    <>
      {file.type === 'file' ? (
        <>
          {file.snapshot?.status === Status.New ? <IconBadgeNew /> : null}
          {file.snapshot?.status === Status.Processing ? (
            <IconBadgeProcessing />
          ) : null}
          {file.snapshot?.status === Status.Error ? <IconBadgeError /> : null}
          {file.isShared ? <IconBadgeShared /> : null}
          {file.snapshot?.entities ? <IconBadgeInsights /> : null}
          {file.snapshot?.mosaic ? <IconBadgeMosaic /> : null}
        </>
      ) : null}
      {file.type === 'folder' ? (
        <>
          {file.isShared ? <IconBadgeShared /> : null}
          {isLoading ? <IconBadgeProcessing /> : null}
        </>
      ) : null}
    </>
  )
}

export default IconBadge
