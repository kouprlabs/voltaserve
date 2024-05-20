import globals from 'globals'
import ts from 'typescript-eslint'

export default [
  ...ts.configs.recommended,
  {
    files: ['src/**/*.{ts}'],
  },
  {
    languageOptions: {
      globals: globals.node,
    },
  },
  {
    ignores: ['*.js'],
  },
]
