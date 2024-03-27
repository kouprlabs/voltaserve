import cx from 'classnames'
import { SnapshotStatus } from '@/client/api/file'
import { FileCommonProps } from '@/types/file'
import IconErrorBadge from './icon-error-badge'
import IconFont from './icon-font'
import IconNewBadge from './icon-new-badge'
import IconProcessingBadge from './icon-processing-badge'
import IconSharedBadge from './icon-shared-badge'
import IconThumbnail from './icon-thumbnail'

type IconFileProps = FileCommonProps

const IconFile = ({ file, scale, viewType }: IconFileProps) => (
  <>
    {file.thumbnail ? (
      <IconThumbnail file={file} scale={scale} />
    ) : (
      <div className={cx('relative')}>
        <IconFont file={file} scale={scale} />
        <div
          className={cx(
            'absolute',
            'flex',
            'flex-row',
            'items-center',
            'gap-[2px]',
            { 'bottom-[-5px]': viewType === 'grid' },
            { 'bottom-[-7px]': viewType === 'list' },
            'right-0',
          )}
        >
          {file.isShared ? <IconSharedBadge /> : null}
          {file.status === SnapshotStatus.New ? <IconNewBadge /> : null}
          {file.status === SnapshotStatus.Processing ? (
            <IconProcessingBadge />
          ) : null}
          {file.status === SnapshotStatus.Error ? <IconErrorBadge /> : null}
        </div>
      </div>
    )}
  </>
)

export default IconFile
