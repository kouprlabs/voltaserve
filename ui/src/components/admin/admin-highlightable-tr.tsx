import { CSSProperties, MouseEvent, ReactNode } from 'react'
import { Tr } from '@chakra-ui/react'
import { useColorModeValue } from '@chakra-ui/system'

export interface AdminHighlightableProps {
  onClick: (event: MouseEvent) => void
  style?: CSSProperties
  children: ReactNode
}

const AdminHighlightableTr = (props: AdminHighlightableProps) => {
  const hoverBg = useColorModeValue('gray.300', 'gray.700')

  return (
    <Tr
      _hover={{
        backgroundColor: hoverBg,
      }}
      style={{ ...props.style, cursor: 'pointer' }}
      onClick={props.onClick}
    >
      {props.children}
    </Tr>
  )
}

export default AdminHighlightableTr
