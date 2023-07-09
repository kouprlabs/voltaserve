import { Circle, Tooltip } from '@chakra-ui/react'
import { HiLanguage } from 'react-icons/hi2'

const OcrSign = () => (
  <Tooltip label="This item has OCR">
    <Circle
      color="MediumSpringGreen"
      bg="white"
      size="23px"
      border="1px solid"
      borderColor="gray.200"
    >
      <HiLanguage fontSize="14px" />
    </Circle>
  </Tooltip>
)

export default OcrSign
