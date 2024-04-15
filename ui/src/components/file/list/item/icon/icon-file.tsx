import cx from 'classnames'
import { SnapshotStatus } from '@/client/api/file'
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
    {file.thumbnail ? (
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
          {file.status === SnapshotStatus.New ? <IconNewBadge /> : null}
          {file.status === SnapshotStatus.Processing ? (
            <IconProcessingBadge />
          ) : null}
          {file.status === SnapshotStatus.Error ? <IconErrorBadge /> : null}
          {file.isShared ? <IconSharedBadge /> : null}
        </div>
      </>
    )}
  </>
)

export default IconFile
