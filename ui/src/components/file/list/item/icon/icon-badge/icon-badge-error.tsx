import { Circle, Tooltip } from '@chakra-ui/react'
import cx from 'classnames'
import { IconClose } from '@/lib/components/icons'

const IconBadgeError = () => (
  <Tooltip label="An error occured while processing this item">
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
      <IconClose className={cx('text-red-600')} />
    </Circle>
  </Tooltip>
)

export default IconBadgeError
