// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

const { variables } = require('./src/lib/variables.cjs')

module.exports = {
  important: true,
  darkMode: 'class',
  content: ['./src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      fontFamily: {
        display: [variables.headingFontFamily],
        body: [variables.bodyFontFamily],
      },
      fontSize: {
        'base': variables.bodyFontSize,
        'heading': variables.headingFontSize,
      },
      borderRadius: {
        'DEFAULT': variables.borderRadius,
        'sm': variables.borderRadiusXs,
        'md': variables.borderRadiusSm,
        'lg': variables.borderRadius,
        'xl': variables.borderRadiusMd,
      },
      spacing: {
        'DEFAULT': variables.spacing,
        '0.5': variables.spacingXs,
        '1': variables.spacingSm,
        '1.5': variables.spacing,
        '2': variables.spacingMd,
        '2.5': variables.spacingLg,
        '3': variables.spacingXl,
        '3.5': variables.spacing2Xl,
      },
    },
  },
  plugins: [],
}
