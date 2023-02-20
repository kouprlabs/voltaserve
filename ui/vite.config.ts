import react from '@vitejs/plugin-react'
import { defineConfig, loadEnv } from 'vite'
import svgr from 'vite-plugin-svgr'
import tsconfigPaths from 'vite-tsconfig-paths'

export default ({ mode }) => {
  process.env = Object.assign(process.env, loadEnv(mode, process.cwd(), ''))
  return defineConfig({
    plugins: [react(), tsconfigPaths(), svgr()],
    server: {
      proxy: {
        '/proxy/api/v1': {
          target: process.env.API_URL,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/proxy\/api/, ''),
        },
        '/proxy/idp/v1': {
          target: process.env.IDP_URL,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/proxy\/idp/, ''),
        },
      },
    },
  })
}
