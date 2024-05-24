import cx from 'classnames'
import { Status } from '@/client/api/snapshot'
import { FileCommonProps } from '@/types/file'
import IconDiverse from './icon-diverse'
import IconErrorBadge from './icon-error-badge'
import IconNewBadge from './icon-new-badge'
import IconProcessingBadge from './icon-processing-badge'
import IconSharedBadge from './icon-shared-badge'
import IconThumbnail from './icon-thumbnail'

type IconFileProps = FileCommonProps

const IconFile = ({ file, scale }: IconFileProps) => (
  <>
    {file.snapshot?.thumbnail ? (
      <IconThumbnail file={file} scale={scale} />
    ) : (
      <>
        <IconDiverse file={file} scale={scale} />
        <div
          className={cx(
            'absolute',
            'flex',
            'flex-row',
            'items-center',
            'gap-[2px]',
            'bottom-[-5px]',
            'right-[-5px]',
          )}
        >
          {file.snapshot?.status === Status.New ? <IconNewBadge /> : null}
          {file.snapshot?.status === Status.Processing ? (
            <IconProcessingBadge />
          ) : null}
          {file.snapshot?.status === Status.Error ? <IconErrorBadge /> : null}
          {file.isShared ? <IconSharedBadge /> : null}
        </div>
      </>
    )}
  </>
)

export default IconFile
