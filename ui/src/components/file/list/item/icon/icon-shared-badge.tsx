import { Circle, Tooltip } from '@chakra-ui/react'
import cx from 'classnames'
import { FiUsers } from 'react-icons/fi'

const IconSharedBadge = () => (
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
      <FiUsers fontSize="12px" />
    </Circle>
  </Tooltip>
)

export default IconSharedBadge
