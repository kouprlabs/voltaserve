import classNames from 'classnames'
import { SnapshotStatus } from '@/client/api/file'
import { CommonItemProps } from '@/types/file'
import ErrorBadge from './error-badge'
import FontIcon from './font-icon'
import NewBadge from './new-badge'
import ProcessingBadge from './processing-badge'
import SharedBadge from './shared-badge'
import Thumbnail from './thumbnail'

type FileIconProps = CommonItemProps

const FileIcon = ({ file, scale, viewType }: FileIconProps) => (
  <>
    {file.thumbnail ? (
      <Thumbnail file={file} scale={scale} />
    ) : (
      <div className={classNames('relative')}>
        <FontIcon file={file} scale={scale} />
        <div
          className={classNames(
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
          {file.isShared ? <SharedBadge /> : null}
          {file.status === SnapshotStatus.New ? <NewBadge /> : null}
          {file.status === SnapshotStatus.Processing ? (
            <ProcessingBadge />
          ) : null}
          {file.status === SnapshotStatus.Error ? <ErrorBadge /> : null}
        </div>
      </div>
    )}
  </>
)

export default FileIcon
