import { useState } from 'react'
import { Image, Skeleton, useColorModeValue, useToken } from '@chakra-ui/react'
import { IconPlay, variables } from '@koupr/ui'
import cx from 'classnames'
import { File, SnapshotStatus } from '@/client/api/file'
import * as fe from '@/helpers/file-extension'
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
  const borderColor = useToken(
    'colors',
    useColorModeValue('gray.300', 'gray.700'),
  )
  return (
    <div className={cx('relative')} style={{ width, height }}>
      <Image
        src={file.thumbnail?.base64}
        style={{
          width: isLoading ? 0 : width,
          height: isLoading ? 0 : height,
          border: '1px solid',
          borderColor,
          borderRadius: variables.borderRadiusSm,
        }}
        className={cx(
          'pointer-events-none',
          'object-cover',
          'border',
          'border-solid',
          {
            'invisible': isLoading,
          },
        )}
        alt={file.name}
        onLoad={() => setIsLoading(false)}
      />
      {isLoading && (
        <Skeleton
          width={width}
          height={height}
          borderRadius={variables.borderRadiusSm}
        />
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
          <IconPlay fontSize="40px" color="white" />
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
