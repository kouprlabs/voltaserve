import { ReactNode } from 'react'
import { Circle } from '@chakra-ui/react'
import { useAppSelector } from '@/store/hook'
import { NavType } from '@/store/ui/nav'
import variables from '@/theme/variables'

type ActiveCircleProps = {
  children: ReactNode
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
