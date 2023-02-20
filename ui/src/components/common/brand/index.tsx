import { useCallback } from 'react'
import { Center, CenterProps, useColorMode } from '@chakra-ui/react'
import { ReactComponent as BrandDarkGlossySvg } from './brand-dark-glossy.svg'
import { ReactComponent as BrandDarkSvg } from './brand-dark.svg'
import { ReactComponent as BrandGlossySvg } from './brand-glossy.svg'
import { ReactComponent as BrandSvg } from './brand.svg'

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
