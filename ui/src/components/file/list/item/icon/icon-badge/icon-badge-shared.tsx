import { Circle, Tooltip } from '@chakra-ui/react'
import cx from 'classnames'
import { IconGroup } from '@/lib'

const IconBadgeShared = () => (
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
      <IconGroup className={cx('text-[12px]')} />
    </Circle>
  </Tooltip>
)

export default IconBadgeShared
