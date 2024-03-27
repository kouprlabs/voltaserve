import { Circle, Tooltip, useToken } from '@chakra-ui/react'
import cx from 'classnames'
import { FiUsers } from 'react-icons/fi'

const IconSharedBadge = () => {
  const borderColor = useToken('colors', 'gray.200')
  return (
    <Tooltip label="This item is shared">
      <Circle
        className={cx(
          'text-orange-600',
          'bg-white',
          'w-[23px]',
          'h-[23px]',
          'border',
          'border-solid',
        )}
        style={{ borderColor }}
      >
        <FiUsers fontSize="12px" />
      </Circle>
    </Tooltip>
  )
}

export default IconSharedBadge
