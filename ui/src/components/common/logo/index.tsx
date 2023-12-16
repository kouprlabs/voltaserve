import { useCallback } from 'react'
import { Center, CenterProps, useColorMode } from '@chakra-ui/react'
import LogoDarkGlossySvg from './logo-dark-glossy.svg?react'
import LogoDarkSvg from './logo-dark.svg?react'
import LogoGlossySvg from './logo-glossy.svg?react'
import LogoSvg from './logo.svg?react'

type LogoProps = CenterProps & {
  isGlossy?: boolean
}

const Logo = ({ isGlossy = false, ...props }: LogoProps) => {
  const { colorMode } = useColorMode()
  const renderSvg = useCallback(() => {
    if (isGlossy) {
      return colorMode === 'dark' ? <LogoDarkGlossySvg /> : <LogoGlossySvg />
    } else {
      return colorMode === 'dark' ? <LogoDarkSvg /> : <LogoSvg />
    }
  }, [colorMode, isGlossy])
  return <Center {...props}>{renderSvg()}</Center>
}

export default Logo
