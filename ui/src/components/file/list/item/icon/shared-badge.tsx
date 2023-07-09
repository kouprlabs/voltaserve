import { Circle, Tooltip } from '@chakra-ui/react'
import { FiUsers } from 'react-icons/fi'

const SharedBadge = () => (
  <Tooltip label="This item is shared">
    <Circle
      color="darkorange"
      bg="white"
      size="23px"
      border="1px solid"
      borderColor="gray.200"
    >
      <FiUsers fontSize="12px" />
    </Circle>
  </Tooltip>
)

export default SharedBadge
