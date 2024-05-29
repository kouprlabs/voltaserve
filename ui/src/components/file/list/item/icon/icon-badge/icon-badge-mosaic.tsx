import { Circle, Tooltip } from '@chakra-ui/react'
import cx from 'classnames'
import { IconModeHeat } from '@/lib/components/icons'

const IconBadgeMosaic = () => (
  <Tooltip label="This item has a mosaic">
    <Circle
      className={cx(
        'text-orange-600',
        'bg-white',
        'w-[23px]',
        'h-[23px]',
        'border',
        'border-gray-200',
      )}
    >
      <IconModeHeat className={cx('text-[12px]')} />
    </Circle>
  </Tooltip>
)

export default IconBadgeMosaic
