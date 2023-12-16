import { useCallback } from 'react'
import { Center, CenterProps, useColorMode } from '@chakra-ui/react'
import BrandDarkGlossySvg from './brand-dark-glossy.svg?react'
import BrandDarkSvg from './brand-dark.svg?react'
import BrandGlossySvg from './brand-glossy.svg?react'
import BrandSvg from './brand.svg?react'

type LogoProps = CenterProps & {
  isGlossy?: boolean
}

const Brand = ({ isGlossy = false, ...props }: LogoProps) => {
  const { colorMode } = useColorMode()
  const renderSvg = useCallback(() => {
    if (isGlossy) {
      return colorMode === 'dark' ? <BrandDarkGlossySvg /> : <BrandGlossySvg />
    } else {
      return colorMode === 'dark' ? <BrandDarkSvg /> : <BrandSvg />
    }
  }, [colorMode, isGlossy])
  return <Center {...props}>{renderSvg()}</Center>
}

export default Brand
