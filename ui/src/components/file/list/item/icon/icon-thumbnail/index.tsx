import { useState } from 'react'
import { Image, Skeleton } from '@chakra-ui/react'
import cx from 'classnames'
import { File, SnapshotStatus } from '@/client/api/file'
import * as fe from '@/helpers/file-extension'
import { IconPlay } from '@/lib'
import IconErrorBadge from '../icon-error-badge'
import IconNewBadge from '../icon-new-badge'
import IconProcessingBadge from '../icon-processing-badge'
import IconSharedBadge from '../icon-shared-badge'
import { getThumbnailHeight, getThumbnailWidth } from './size'

export type IconThumbnailProps = {
  file: File
  scale: number
}

const IconThumbnail = ({ file, scale }: IconThumbnailProps) => {
  const { original, isShared, status } = file
  const width = getThumbnailWidth(file, scale)
  const height = getThumbnailHeight(file, scale)
  const [isLoading, setIsLoading] = useState(true)

  return (
    <div className={cx('relative')} style={{ width, height }}>
      <Image
        src={file.thumbnail?.base64}
        style={{
          width: isLoading ? 0 : width,
          height: isLoading ? 0 : height,
        }}
        className={cx(
          'pointer-events-none',
          'object-cover',
          'border',
          'border-solid',
          { 'invisible': isLoading },
          'border',
          'border-gray-300',
          'dark:border-gray-700',
          'rounded-md',
        )}
        alt={file.name}
        onLoad={() => setIsLoading(false)}
      />
      {isLoading && (
        <Skeleton className={cx('rounded-md')} style={{ width, height }} />
      )}
      {fe.isVideo(original?.extension) && (
        <div
          className={cx(
            'absolute',
            'top-0',
            'left-0',
            'opacity-50',
            'flex',
            'items-center',
            'justify-center',
          )}
          style={{ width, height }}
        >
          <IconPlay fontSize="40px" className={cx('text-white')} />
        </div>
      )}
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
        {isShared ? <IconSharedBadge /> : null}
        {status === SnapshotStatus.New ? <IconNewBadge /> : null}
        {status === SnapshotStatus.Processing ? <IconProcessingBadge /> : null}
        {status === SnapshotStatus.Error ? <IconErrorBadge /> : null}
      </div>
    </div>
  )
}

export default IconThumbnail
