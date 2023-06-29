import { ReactNode } from 'react'
import { Circle } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { useAppSelector } from '@/store/hook'
import { NavType } from '@/store/ui/nav'

type ActiveCircleProps = {
  children?: ReactNode
}

const ActiveCircle = ({ children }: ActiveCircleProps) => {
  const activeNav = useAppSelector((state) => state.ui.nav.active)
  return (
    <Circle
      size="50px"
      bg={activeNav === NavType.Account ? variables.gradiant : 'transparent'}
    >
      {children}
    </Circle>
  )
}

export default ActiveCircle
