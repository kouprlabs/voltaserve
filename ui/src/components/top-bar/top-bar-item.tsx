import { Link } from 'react-router-dom'
import { Link as ChakraLink, useColorModeValue } from '@chakra-ui/react'
import { variables } from '@koupr/ui'

type TopBarItemProps = {
  title: string
  href: string
  isActive: boolean
}

const TopBarItem = ({ title, href, isActive }: TopBarItemProps) => {
  const activeTextColor = useColorModeValue('white', 'gray.800')
  const activeBgColor = useColorModeValue('black', 'white')
  return (
    <ChakraLink
      as={Link}
      to={href}
      padding="0 20px 0 20px"
      bg={isActive ? activeBgColor : 'transparent'}
      color={isActive ? activeTextColor : ''}
      fontWeight="semibold"
      h="40px"
      lineHeight="40px"
      borderRadius={variables.borderRadiusMd}
      variant="no-underline"
      className="opacity-100 hover:opacity-80"
    >
      {title}
    </ChakraLink>
  )
}

export default TopBarItem
