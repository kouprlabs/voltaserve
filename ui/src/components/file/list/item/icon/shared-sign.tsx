import { Circle } from '@chakra-ui/react'
import { FiUsers } from 'react-icons/fi'

const SharedSign = () => (
  <Circle
    color="darkorange"
    bg="white"
    size="23px"
    border="1px solid"
    borderColor="gray.200"
  >
    <FiUsers fontSize="12px" />
  </Circle>
)

export default SharedSign
