import { useCallback } from 'react'
import { Center, CenterProps, useColorMode } from '@chakra-ui/react'
import { ReactComponent as LogoDarkGlossySvg } from './logo-dark-glossy.svg'
import { ReactComponent as LogoDarkSvg } from './logo-dark.svg'
import { ReactComponent as LogoGlossySvg } from './logo-glossy.svg'
import { ReactComponent as LogoSvg } from './logo.svg'

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
