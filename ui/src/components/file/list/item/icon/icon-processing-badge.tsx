import { Circle, Spinner, Tooltip } from '@chakra-ui/react'
import cx from 'classnames'

const IconProcessingBadge = () => (
  <Tooltip label="Processing in progress">
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
      <Spinner
        size="sm"
        thickness="4px"
        className={cx('text-gray-400', 'dark:text-gray-500')}
      />
    </Circle>
  </Tooltip>
)

export default IconProcessingBadge
