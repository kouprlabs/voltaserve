import { ReactElement, useContext, useEffect, useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { Center, Stack, Text, useColorModeValue } from '@chakra-ui/react'
import { DrawerContext } from '@/components/common/drawer'
import variables from '@/theme/variables'

type ItemProps = {
  icon: ReactElement
  href: string
  primaryText: string
  secondaryText: string
  isActive?: boolean
}

const Item = ({ icon, href, primaryText, secondaryText }: ItemProps) => {
  const location = useLocation()
  const [isActive, setIsActive] = useState<boolean>()
  const { isCollapsed } = useContext(DrawerContext)
  const bgPressedColor = useColorModeValue('gray.200', 'gray.700')
  const bgHoverColor = useColorModeValue('gray.100', 'gray.600')
  const bgActiveColor = useColorModeValue('black', 'white')
  const textActiveColor = useColorModeValue('white', 'gray.800')

  useEffect(() => {
    if (
      (href === '/' && location.pathname === '/') ||
      (href !== '/' && location.pathname.startsWith(href))
    ) {
      setIsActive(true)
    } else {
      setIsActive(false)
    }
  }, [location.pathname, href])

  return (
    <Link
      to={href}
      style={{ width: '100%' }}
      title={isCollapsed ? `${primaryText}: ${secondaryText}` : secondaryText}
    >
      <Stack
        direction="row"
        spacing={variables.spacing}
        _hover={{ bg: isActive ? '' : bgHoverColor }}
        _active={{ bg: isActive ? '' : bgPressedColor }}
        bg={isActive ? bgActiveColor : ''}
        borderRadius={variables.borderRadiusSm}
        p={variables.spacing}
        minW="50px"
        h="50px"
      >
        <Center flexShrink={0} color={isActive ? textActiveColor : ''}>
          {icon}
        </Center>
        {!isCollapsed && (
          <Text color={isActive ? textActiveColor : ''}>{primaryText}</Text>
        )}
      </Stack>
    </Link>
  )
}

export default Item
