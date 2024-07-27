// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import eslintPluginPrettierRecommended from 'eslint-plugin-prettier/recommended'
import react from 'eslint-plugin-react'
import globals from 'globals'
import ts from 'typescript-eslint'

export default [
  ...ts.configs.recommended,
  eslintPluginPrettierRecommended,
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
      'prettier/prettier': 'error',
    },
  },
  {
    ignores: ['*.cjs', 'dist'],
  },
]
