import alias from '@rollup/plugin-alias'
import commonjs from '@rollup/plugin-commonjs'
import image from '@rollup/plugin-image'
import resolve from '@rollup/plugin-node-resolve'
import strip from '@rollup/plugin-strip'
import typescript from '@rollup/plugin-typescript'
import { createFilter } from '@rollup/pluginutils'
import { createRequire } from 'module'
import path from 'path'
import { dts } from 'rollup-plugin-dts'
import peerDepsExternal from 'rollup-plugin-peer-deps-external'
import postcss from 'rollup-plugin-postcss'

const require = createRequire(import.meta.url)
const pkg = require('./package.json')

const injectCssImportPlugin = () => {
  return {
    name: 'inject-css-import',
    generateBundle(options, bundle) {
      const filter = createFilter('**/*.js')
      for (const file of Object.keys(bundle)) {
        if (filter(file)) {
          const code = bundle[file].code
          // Determine the CSS file name by replacing .js with .css in each JS file name
          const cssFileName = path.basename(file).replace('.js', '.css')
          // Modify the bundle to include the CSS import statement
          bundle[file].code = `import './${cssFileName}';\n${code}`
        }
      }
    },
  }
}

export default [
  {
    input: 'src/index.tsx',
    output: [
      {
        file: pkg.main,
        format: 'cjs',
        sourcemap: true,
      },
      {
        file: pkg.module,
        format: 'es',
        sourcemap: true,
      },
    ],
    external: [
      '@chakra-ui/anatomy',
      '@chakra-ui/react',
      '@chakra-ui/theme-tools',
      '@dnd-kit/core',
      '@dnd-kit/modifiers',
      '@dnd-kit/sortable',
      '@emotion/css',
      '@emotion/react',
      '@emotion/styled',
      '@google/model-viewer',
      '@koupr/ui',
      '@nivo/core',
      '@nivo/pie',
      '@reduxjs/toolkit',
      'chakra-react-select',
      'classnames',
      'formik',
      'framer-motion',
      'hashids',
      'jose',
      'js-base64',
      'react',
      'react-dom',
      'react-dropzone',
      'react-helmet-async',
      'react-hotkeys-hook',
      'react-redux',
      'react-router-dom',
      'react-use',
      'redux',
      'semver',
      'swr',
      'uuid',
      'yup',
    ],
    plugins: [
      peerDepsExternal(),
      alias({
        entries: [{ find: '@', replacement: path.resolve(__dirname, 'src') }],
      }),
      resolve(),
      commonjs(),
      typescript({
        tsconfig: 'tsconfig.rollup.json',
        include: ['src/**/*.{ts,tsx}'],
        exclude: ['src/main.tsx'],
      }),
      image(),
      postcss({
        extract: true,
      }),
      injectCssImportPlugin(),
      strip(),
    ],
  },
  {
    input: 'dist/index.d.ts',
    output: [{ file: 'dist/types.d.ts', format: 'es' }],
    plugins: [
      dts({
        tsconfig: 'tsconfig.rollup.json',
      }),
    ],
    external: [/\.css$/],
  },
]
