import classNames from 'classnames'
import { FcFolder } from 'react-icons/fc'
import { FileCommonProps } from '@/types/file'
import IconProcessingBadge from './icon-processing-badge'
import IconSharedBadge from './icon-shared-badge'

type IconFolderProps = {
  isLoading?: boolean
} & FileCommonProps

const ICON_FONT_SIZE = 92

const IconFolder = ({ file, scale, viewType, isLoading }: IconFolderProps) => (
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
      {file.isShared ? <IconSharedBadge /> : null}
      {isLoading ? <IconProcessingBadge /> : null}
    </div>
  </div>
)

export default IconFolder
