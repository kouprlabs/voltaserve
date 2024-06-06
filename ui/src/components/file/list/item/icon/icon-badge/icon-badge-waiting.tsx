import { Circle, Tooltip } from '@chakra-ui/react'
import cx from 'classnames'
import { IconHourglass } from '@/lib/components/icons'

const IconBadgeWaiting = () => (
  <Tooltip label="Waiting for processing">
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
      <IconHourglass />
    </Circle>
  </Tooltip>
)

export default IconBadgeWaiting
