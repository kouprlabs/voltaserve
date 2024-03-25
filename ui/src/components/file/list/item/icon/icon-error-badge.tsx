import { Circle, Tooltip } from '@chakra-ui/react'
import { IoClose } from 'react-icons/io5'

const IconErrorBadge = () => (
  <Tooltip label="An error occured while processing this item">
    <Circle
      color="darkorange"
      bg="white"
      size="23px"
      border="1px solid"
      borderColor="gray.200"
    >
      <IoClose fontSize="14px" color="red" />
    </Circle>
  </Tooltip>
)

export default IconErrorBadge
