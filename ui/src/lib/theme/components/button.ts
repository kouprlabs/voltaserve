import { mode, StyleFunctionProps } from '@chakra-ui/theme-tools'
import variables from '../../variables'

const button = {
  baseStyle: {
    borderRadius: variables.borderRadiusMd,
    fontWeight: variables.bodyFontWeight,
  },
  sizes: {
    md: {
      fontSize: variables.bodyFontSize,
    },
    xs: {
      fontSize: '12px',
    },
  },
  variants: {
    'solid-gray': (props: StyleFunctionProps) => ({
      bg: mode('gray.100', 'gray.700')(props),
      _hover: {
        bg: mode('gray.200', 'gray.600')(props),
      },
    }),
  },
}

export default button
