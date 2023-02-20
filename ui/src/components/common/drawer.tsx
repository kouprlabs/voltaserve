import { createContext, ReactNode, useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import {
  Box,
  Center,
  Flex,
  Stack,
  useColorModeValue,
  VStack,
} from '@chakra-ui/react'
import { BsChevronBarLeft, BsChevronBarRight } from 'react-icons/bs'
import Logo from '@/components/common/logo'
import variables from '@/theme/variables'

type DrawerContextType = {
  isCollapsed: boolean | undefined
  isTouched: boolean
}

export const DrawerContext = createContext<DrawerContextType>({
  isCollapsed: undefined,
  isTouched: false,
})

type DrawerProps = {
  children?: ReactNode
  localStorageNamespace: string
}

const Drawer = ({ children, localStorageNamespace }: DrawerProps) => {
  const [isCollapsed, setIsCollapsed] = useState<boolean | undefined>(undefined)
  const [isTouched, setIsTouched] = useState(false)
  const borderColor = useColorModeValue('gray.200', 'gray.700')
  const activeColor = useColorModeValue('gray.200', 'gray.700')
  const hoverColor = useColorModeValue('gray.100', 'gray.600')
  const localStorageCollapsedKey = useMemo(
    () => `voltaserve_${localStorageNamespace}_drawer_collapsed`,
    [localStorageNamespace]
  )

  useEffect(() => {
    let collapse = false
    if (typeof localStorage !== 'undefined') {
      const value = localStorage.getItem(localStorageCollapsedKey)
      if (value) {
        collapse = JSON.parse(value)
      } else {
        localStorage.setItem(localStorageCollapsedKey, JSON.stringify(false))
      }
    }
    if (collapse) {
      setIsCollapsed(true)
    } else {
      setIsCollapsed(false)
    }
  }, [localStorageCollapsedKey, setIsCollapsed])

  if (isCollapsed === undefined) {
    return null
  }

  return (
    <DrawerContext.Provider
      value={{
        isCollapsed,
        isTouched,
      }}
    >
      <Stack
        direction="column"
        h="100%"
        borderRight="1px solid"
        borderRightColor={borderColor}
        flexShrink={0}
        spacing={0}
      >
        <Center h="80px" alignItems="center" justifyContent="center">
          <Link to="/">
            <Flex h="40px">
              <Center w="40px" h="40px">
                <Logo className="w-4" />
              </Center>
            </Flex>
          </Link>
        </Center>
        <VStack
          direction="column"
          spacing={variables.spacingXs}
          p={`0 ${variables.spacing} ${variables.spacing} ${variables.spacing}`}
        >
          {children}
        </VStack>
        <Box flexGrow={1} />
        <Stack
          h="50px"
          w="100%"
          _hover={{ bg: hoverColor }}
          _active={{ bg: activeColor }}
          direction="row"
          alignItems="center"
          justifyContent={isCollapsed ? 'center' : 'flex-end'}
          px={isCollapsed ? 0 : variables.spacing}
          cursor="pointer"
          onClick={() => {
            setIsCollapsed(!isCollapsed)
            setIsTouched(true)
            localStorage.setItem(
              localStorageCollapsedKey,
              JSON.stringify(!isCollapsed)
            )
          }}
        >
          {isCollapsed ? (
            <BsChevronBarRight fontSize="16px" />
          ) : (
            <BsChevronBarLeft fontSize="16px" />
          )}
        </Stack>
      </Stack>
    </DrawerContext.Provider>
  )
}

export default Drawer
