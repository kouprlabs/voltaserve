import { Circle, Spinner, Tooltip, useColorModeValue } from '@chakra-ui/react'

const IconProcessingBadge = () => {
  const spinnerColor = useColorModeValue('gray.400', 'gray.500')
  return (
    <Tooltip label="Processing in progress">
      <Circle bg="white" size="23px" border="1px solid" borderColor="gray.200">
        <Spinner size="sm" thickness="4px" color={spinnerColor} />
      </Circle>
    </Tooltip>
  )
}

export default IconProcessingBadge
