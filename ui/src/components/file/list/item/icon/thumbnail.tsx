import { useState } from 'react'
import { Image, Skeleton, useColorModeValue, useToken } from '@chakra-ui/react'
import { IconPlay, variables } from '@koupr/ui'
import classNames from 'classnames'
import { File, SnapshotStatus } from '@/client/api/file'
import { getSizeWithAspectRatio } from '@/helpers/aspect-ratio'
import * as fe from '@/helpers/file-extension'
import ErrorBadge from './error-badge'
import NewBadge from './new-badge'
import ProcessingBadge from './processing-badge'
import SharedBadge from './shared-badge'

const MAX_WIDTH = 130
const MAX_HEIGHT = 130

export function getThumbnailWidth(file: File, scale: number): string {
  if (file.thumbnail) {
    const { width } = getSizeWithAspectRatio(
      file.thumbnail.width,
      file.thumbnail.height,
      MAX_WIDTH,
      MAX_HEIGHT,
    )
    return `${width * scale}px`
  } else {
    return `${MAX_WIDTH * scale}px`
  }
}

export function getThumbnailHeight(file: File, scale: number): string {
  if (file.thumbnail) {
    const { height } = getSizeWithAspectRatio(
      file.thumbnail.width,
      file.thumbnail.height,
      MAX_WIDTH,
      MAX_HEIGHT,
    )
    return `${height * scale}px`
  } else {
    return `${MAX_HEIGHT * scale}px`
  }
}

type ThumbnailProps = {
  file: File
  scale: number
}

const Thumbnail = ({ file, scale }: ThumbnailProps) => {
  const { original, isShared, status } = file
  const width = getThumbnailWidth(file, scale)
  const height = getThumbnailHeight(file, scale)
  const [isLoading, setIsLoading] = useState(true)
  const borderColor = useToken(
    'colors',
    useColorModeValue('gray.300', 'gray.700'),
  )
  return (
    <div className={classNames('relative')} style={{ width, height }}>
      <Image
        src={file.thumbnail?.base64}
        width={isLoading ? 0 : width}
        height={isLoading ? 0 : height}
        style={{
          objectFit: 'cover',
          width: isLoading ? 0 : width,
          height: isLoading ? 0 : height,
          border: '1px solid',
          borderColor,
          borderRadius: variables.borderRadiusSm,
          visibility: isLoading ? 'hidden' : undefined,
        }}
        pointerEvents="none"
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
          className={classNames(
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
        className={classNames(
          'absolute',
          'flex',
          'flex-row',
          'items-center',
          'gap-[2px]',
          'bottom-[-5px]',
          'right-[-5px]',
        )}
      >
        {isShared ? <SharedBadge /> : null}
        {status === SnapshotStatus.New ? <NewBadge /> : null}
        {status === SnapshotStatus.Processing ? <ProcessingBadge /> : null}
        {status === SnapshotStatus.Error ? <ErrorBadge /> : null}
      </div>
    </div>
  )
}

export default Thumbnail
