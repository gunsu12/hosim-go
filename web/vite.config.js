import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import tailwindcss from '@tailwindcss/vite'

import { exec } from 'node:child_process'
import path from 'path'

function openInAntigravityPlugin() {
  return {
    name: 'open-in-antigravity',
    configureServer(server) {
      server.middlewares.use('/__open-in-editor', (req, res, next) => {
        try {
          const fullUrl = new URL(req.url, 'http://localhost')
          const file = fullUrl.searchParams.get('file')
          if (file) {
            exec(`antigravity-ide.cmd -r -g "${file}"`, (err) => {
              if (err) {
                console.error('[Inspector] Gagal membuka Antigravity IDE:', err)
              }
            })
            res.statusCode = 200
            res.end('OK')
            return
          }
        } catch (e) {
          console.error('[Inspector] Error handling open in editor:', e)
        }
        next()
      })
    },
  }
}

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    openInAntigravityPlugin(),
    tailwindcss(),
    svelte({
      inspector: {
        toggleKeyCombo: 'control-shift',
        holdMode: true,
        showToggleButton: 'always',
        toggleButtonPos: 'bottom-right',
      },
    }),
  ],
  resolve: {
    alias: {
      '$lib': path.resolve(import.meta.dirname, './src/lib'),
    },
  },
  server: {
    port: 5125,
    proxy: {
      '/api': {
        target: 'http://localhost:8910',
        changeOrigin: true,
      },
    },
  },
})
