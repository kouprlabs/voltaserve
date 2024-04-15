import cx from 'classnames'
import { FileCommonProps } from '@/types/file'
import { computeScale } from '../scale'
import IconFile from './icon-file'
import IconFolder from './icon-folder'

export type ItemIconProps = {
  isLoading?: boolean
} & FileCommonProps

const ItemIcon = ({ file, scale, viewType, isLoading }: ItemIconProps) => (
  <>
    <div
      className={cx('z-0', 'text-gray-500', 'dark:text-gray-300', 'relative')}
    >
      {file.type === 'file' ? (
        <IconFile
          file={file}
          scale={computeScale(scale, viewType)}
          viewType={viewType}
        />
      ) : file.type === 'folder' ? (
        <IconFolder
          file={file}
          scale={computeScale(scale, viewType)}
          viewType={viewType}
          isLoading={isLoading}
        />
      ) : null}
    </div>
  </>
)

export default ItemIcon
