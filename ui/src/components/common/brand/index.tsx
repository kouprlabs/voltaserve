import { useCallback } from 'react'
import { CenterProps, useColorMode } from '@chakra-ui/react'
import classNames from 'classnames'
import BrandDarkGlossySvg from './brand-dark-glossy.svg?react'
import BrandDarkSvg from './brand-dark.svg?react'
import BrandGlossySvg from './brand-glossy.svg?react'
import BrandSvg from './brand.svg?react'

type LogoProps = CenterProps & {
  isGlossy?: boolean
}

const Brand = ({ isGlossy = false }: LogoProps) => {
  const { colorMode } = useColorMode()
  const renderSvg = useCallback(() => {
    if (isGlossy) {
      return colorMode === 'dark' ? <BrandDarkGlossySvg /> : <BrandGlossySvg />
    } else {
      return colorMode === 'dark' ? <BrandDarkSvg /> : <BrandSvg />
    }
  }, [colorMode, isGlossy])
  return (
    <div className={classNames('flex', 'items-center', 'justify-center')}>
      {renderSvg()}
    </div>
  )
}

export default Brand
