import { Circle, Tooltip } from '@chakra-ui/react'
import cx from 'classnames'
import { BsHourglassSplit } from 'react-icons/bs'

const IconNewBadge = () => (
  <Tooltip label="Waiting for processing">
    <Circle
      className={cx(
        'text-purple-600',
        'bg-white',
        'w-[23px]',
        'h-[23px]',
        'border',
        'border-gray-200',
      )}
    >
      <BsHourglassSplit className={cx('text-[14px]')} />
    </Circle>
  </Tooltip>
)

export default IconNewBadge
