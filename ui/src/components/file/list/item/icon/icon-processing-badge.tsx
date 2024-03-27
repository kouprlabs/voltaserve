import {
  Circle,
  Spinner,
  Tooltip,
  useColorModeValue,
  useToken,
} from '@chakra-ui/react'
import classNames from 'classnames'

const IconProcessingBadge = () => {
  const spinnerColor = useToken(
    'colors',
    useColorModeValue('gray.400', 'gray.500'),
  )
  const borderColor = useToken('colors', 'gray.200')
  return (
    <Tooltip label="Processing in progress">
      <Circle
        className={classNames(
          'text-purple-600',
          'bg-white',
          'w-[23px]',
          'h-[23px]',
          'border',
          'border-solid',
        )}
        style={{ borderColor }}
      >
        <Spinner size="sm" thickness="4px" style={{ color: spinnerColor }} />
      </Circle>
    </Tooltip>
  )
}

export default IconProcessingBadge
