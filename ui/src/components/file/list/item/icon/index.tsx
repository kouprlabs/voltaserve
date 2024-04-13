import cx from 'classnames'
import { FileCommonProps } from '@/types/file'
import IconFile from './icon-file'
import IconFolder from './icon-folder'

export type ItemIconProps = {
  isLoading?: boolean
} & FileCommonProps

const ItemIcon = ({ file, scale, viewType, isLoading }: ItemIconProps) => (
  <>
    {file.type === 'file' ? (
      <div className={cx('z-0', 'text-gray-500', 'dark:text-gray-300')}>
        <IconFile file={file} scale={scale} viewType={viewType} />
      </div>
    ) : file.type === 'folder' ? (
      <div className={cx('z-0', 'text-gray-500', 'dark:text-gray-300')}>
        <IconFolder
          file={file}
          scale={scale}
          viewType={viewType}
          isLoading={isLoading}
        />
      </div>
    ) : null}
  </>
)

export default ItemIcon
