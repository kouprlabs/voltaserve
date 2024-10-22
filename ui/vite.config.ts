// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import react from '@vitejs/plugin-react'
import { defineConfig, loadEnv } from 'vite'
import svgr from 'vite-plugin-svgr'
import tsconfigPaths from 'vite-tsconfig-paths'

const config = ({ mode }) => {
  process.env = Object.assign(process.env, loadEnv(mode, process.cwd(), ''))
  return defineConfig({
    plugins: [react(), tsconfigPaths(), svgr()],
    server: {
      port: 3000,
      proxy: {
        '/proxy/api/v3': {
          target: process.env.API_URL,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/proxy\/api/, ''),
        },
        '/proxy/idp/v3': {
          target: process.env.IDP_URL,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/proxy\/idp/, ''),
        },
        '/proxy/console/v3': {
          target: process.env.CONSOLE_URL,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/proxy\/console/, ''),
        },
      },
    },
  })
}

export default config
