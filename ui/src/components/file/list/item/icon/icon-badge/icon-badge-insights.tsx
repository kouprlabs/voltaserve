import { Circle, Tooltip } from '@chakra-ui/react'
import cx from 'classnames'
import { getClassName } from '@/lib/components/icons'

const IconBadgeInsights = () => (
  <Tooltip label="This item is shared">
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
      <span className={getClassName({ className: 'visibility' })}></span>
    </Circle>
  </Tooltip>
)

export default IconBadgeInsights
