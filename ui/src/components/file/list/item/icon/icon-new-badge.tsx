import { Circle, Tooltip } from '@chakra-ui/react'
import { BsHourglassSplit } from 'react-icons/bs'

const IconNewBadge = () => (
  <Tooltip label="Waiting for processing">
    <Circle
      color="#9B59B6"
      bg="white"
      size="23px"
      border="1px solid"
      borderColor="gray.200"
    >
      <BsHourglassSplit fontSize="14px" />
    </Circle>
  </Tooltip>
)

export default IconNewBadge
