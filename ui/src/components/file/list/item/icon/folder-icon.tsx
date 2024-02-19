import classNames from 'classnames'
import { FcFolder } from 'react-icons/fc'
import { CommonItemProps } from '@/types/file'
import ProcessingBadge from './processing-badge'
import SharedBadge from './shared-badge'

type FolderIconProps = {
  isLoading?: boolean
} & CommonItemProps

const ICON_FONT_SIZE = 92

const FolderIcon = ({ file, scale, viewType, isLoading }: FolderIconProps) => (
  <div className={classNames('relative')}>
    <FcFolder fontSize={`${ICON_FONT_SIZE * scale}px`} />
    <div
      className={classNames(
        'absolute',
        'flex',
        'flex-row',
        'items-center',
        'gap-[2px]',
        { 'bottom-[7px]': viewType === 'grid' },
        { 'right-[2px]': viewType === 'grid' },
        { 'bottom-[0px]': viewType === 'list' },
        { 'right-[-2px]': viewType === 'list' },
      )}
    >
      {file.isShared ? <SharedBadge /> : null}
      {isLoading ? <ProcessingBadge /> : null}
    </div>
  </div>
)

export default FolderIcon
