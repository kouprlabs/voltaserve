import react from 'eslint-plugin-react'
import globals from 'globals'
import ts from 'typescript-eslint'

export default [
  ...ts.configs.recommended,
  {
    files: ['src/**/*.{ts,tsx}'],
  },
  {
    plugins: {
      react,
    },
  },
  {
    languageOptions: {
      globals: globals.browser,
    },
  },
  {
    rules: {
      'react/react-in-jsx-scope': 'off',
    },
  },
  {
    ignores: ['*.cjs'],
  },
]
