import cx from 'classnames'
import { Status } from '@/client/api/snapshot'
import { FileCommonProps } from '@/types/file'
import IconBadgeError from './icon-badge/icon-badge-error'
import IconBadgeNew from './icon-badge/icon-badge-new'
import IconBadgeProcessing from './icon-badge/icon-badge-processing'
import IconBadgeShared from './icon-badge/icon-badge-shared'
import IconDiverse from './icon-diverse'
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
          {file.snapshot?.status === Status.New ? <IconBadgeNew /> : null}
          {file.snapshot?.status === Status.Processing ? (
            <IconBadgeProcessing />
          ) : null}
          {file.snapshot?.status === Status.Error ? <IconBadgeError /> : null}
          {file.isShared ? <IconBadgeShared /> : null}
        </div>
      </>
    )}
  </>
)

export default IconFile
