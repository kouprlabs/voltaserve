import { Circle } from '@chakra-ui/react'
import { FiUsers } from 'react-icons/fi'

type FileListItemSharedSignProps = {
  top?: string
  right?: string
  bottom?: string
  left?: string
}

const FileListItemSharedSign = ({
  top,
  right,
  bottom,
  left,
}: FileListItemSharedSignProps) => (
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

export default FileListItemSharedSign
