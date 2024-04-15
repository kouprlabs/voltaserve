const { variables } = require('./src/lib/variables.cjs')

module.exports = {
  important: true,
  darkMode: 'class',
  content: ['./src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      fontFamily: {
        display: ['IBM Plex Serif'],
        body: ['IBM Plex Sans'],
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
