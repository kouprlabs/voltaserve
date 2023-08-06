import { Circle, Tooltip } from '@chakra-ui/react'
import { Spinner } from '@koupr/ui'

const ProcessingBadge = () => (
  <Tooltip label="Processing in progress">
    <Circle bg="white" size="23px" border="1px solid" borderColor="gray.200">
      <Spinner />
    </Circle>
  </Tooltip>
)

export default ProcessingBadge
