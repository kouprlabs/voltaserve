const link = {
  baseStyle: {
    textDecoration: 'underline',
    _active: {
      boxShadow: 'none',
    },
    _focus: {
      boxShadow: 'none',
    },
  },
  variants: {
    'no-underline': {
      textDecoration: 'none',
      _hover: {
        textDecoration: 'none',
      },
    },
  },
}

export default link
