import { Circle, Tooltip } from '@chakra-ui/react'
import cx from 'classnames'
import { IconVisibility } from '@/lib/components/icons'

const IconBadgeInsights = () => (
  <Tooltip label="This item has insights">
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
      <IconVisibility className={cx('text-[12px]')} />
    </Circle>
  </Tooltip>
)

export default IconBadgeInsights
