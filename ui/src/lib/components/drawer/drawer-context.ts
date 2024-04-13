import { createContext } from 'react'

type DrawerContextType = {
  isCollapsed: boolean | undefined
  isTouched: boolean
}

export const DrawerContext = createContext<DrawerContextType>({
  isCollapsed: undefined,
  isTouched: false,
})
