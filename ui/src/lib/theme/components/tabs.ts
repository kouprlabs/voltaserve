import { mode, StyleFunctionProps } from '@chakra-ui/theme-tools'
import { variables } from '../../variables'

const tab = {
  variants: {
    'solid-rounded': (props: StyleFunctionProps) => ({
      tab: {
        fontSize: variables.bodyFontSize,
        _focus: {
          boxShadow: 'none',
        },
        _selected: {
          bg: mode('black', 'white')(props),
        },
      },
      tabpanel: {
        p: '60px 0 0 0',
      },
    }),
    'line': {
      tab: {
        fontSize: variables.bodyFontSize,
        _focus: {
          boxShadow: 'none',
        },
      },
    },
  },
}

export default tab
