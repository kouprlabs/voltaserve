import { Circle, Tooltip } from '@chakra-ui/react'
import cx from 'classnames'
import { IoClose } from 'react-icons/io5'

const IconErrorBadge = () => (
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
      <IoClose className={cx('text-red-600', 'text-[14px]')} />
    </Circle>
  </Tooltip>
)

export default IconErrorBadge
