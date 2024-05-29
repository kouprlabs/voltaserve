import cx from 'classnames'
import { FileCommonProps } from '@/types/file'
import IconBadgeProcessing from '../icon-badge/icon-badge-processing'
import IconBadgeShared from '../icon-badge/icon-badge-shared'
import FolderSvg from './assets/icon-folder.svg'

type IconFolderProps = {
  isLoading?: boolean
} & FileCommonProps

const MIN_WIDTH = 45
const MIN_HEIGHT = 36.05
const BASE_WIDTH = 67
const BASE_HEIGHT = 53.68

const IconFolder = ({ file, scale, isLoading }: IconFolderProps) => {
  const width = Math.max(MIN_WIDTH, BASE_WIDTH * scale)
  const height = Math.max(MIN_HEIGHT, BASE_HEIGHT * scale)

  return (
    <>
      <img
        src={FolderSvg}
        style={{ width: `${width}px`, height: `${height}px` }}
      />
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
        {file.isShared ? <IconBadgeShared /> : null}
        {isLoading ? <IconBadgeProcessing /> : null}
      </div>
    </>
  )
}

export default IconFolder
