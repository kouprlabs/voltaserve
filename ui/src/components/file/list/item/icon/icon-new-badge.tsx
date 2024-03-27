import { Circle, Tooltip, useToken } from '@chakra-ui/react'
import cx from 'classnames'
import { BsHourglassSplit } from 'react-icons/bs'

const IconNewBadge = () => {
  const borderColor = useToken('colors', 'gray.200')
  return (
    <Tooltip label="Waiting for processing">
      <Circle
        className={cx(
          'text-purple-600',
          'bg-white',
          'w-[23px]',
          'h-[23px]',
          'border',
          'border-solid',
        )}
        style={{ borderColor }}
      >
        <BsHourglassSplit className={cx('text-[14px]')} />
      </Circle>
    </Tooltip>
  )
}

export default IconNewBadge
