import { useState } from 'react'
import { Image, Skeleton } from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import { Status } from '@/client/api/snapshot'
import * as fe from '@/helpers/file-extension'
import { IconPlayArrow } from '@/lib'
import IconBadgeError from '../icon-badge/icon-badge-error'
import IconBadgeNew from '../icon-badge/icon-badge-new'
import IconBadgeProcessing from '../icon-badge/icon-badge-processing'
import IconBadgeShared from '../icon-badge/icon-badge-shared'
import { getThumbnailHeight, getThumbnailWidth } from './size'

export type IconThumbnailProps = {
  file: File
  scale: number
}

const IconThumbnail = ({ file, scale }: IconThumbnailProps) => {
  const { isShared } = file
  const { original, status } = file.snapshot || {}
  const width = getThumbnailWidth(file, scale)
  const height = getThumbnailHeight(file, scale)
  const [isLoading, setIsLoading] = useState(true)
  return (
    <>
      <Image
        src={file.snapshot?.thumbnail?.base64}
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
      {isLoading ? (
        <Skeleton className={cx('rounded-md')} style={{ width, height }} />
      ) : null}
      {fe.isVideo(original?.extension) ? (
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
          <IconPlayArrow
            className={cx('text-white', 'text-[40px]')}
            filled={true}
          />
        </div>
      ) : null}
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
        {isShared ? <IconBadgeShared /> : null}
        {status === Status.New ? <IconBadgeNew /> : null}
        {status === Status.Processing ? <IconBadgeProcessing /> : null}
        {status === Status.Error ? <IconBadgeError /> : null}
      </div>
    </>
  )
}

export default IconThumbnail
