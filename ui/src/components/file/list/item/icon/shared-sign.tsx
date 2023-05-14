import { Circle } from '@chakra-ui/react'
import { FiUsers } from 'react-icons/fi'

type SharedSignProps = {
  top?: string
  right?: string
  bottom?: string
  left?: string
}

const SharedSign = ({ top, right, bottom, left }: SharedSignProps) => (
  <Circle
    position="absolute"
    top={top}
    right={right}
    bottom={bottom}
    left={left}
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
