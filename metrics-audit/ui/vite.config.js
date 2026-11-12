import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

const port = Number(process.env.METRICS_DEV_PORT || 8013)
const origin = process.env.METRICS_DEV_ORIGIN || `http://172.16.1.3:${port}`
const panelOrigin = process.env.METRICS_PANEL_ORIGIN || 'http://218.23.2.48:9090'

export default defineConfig(({ command }) => ({
  base: command === 'serve' ? `${origin}/` : './',
  plugins: [vue(), {
    name: 'local-panel-cors',
    configureServer(server) {
      server.middlewares.use((request, response, next) => {
        response.setHeader('Access-Control-Allow-Origin', panelOrigin)
        response.setHeader('Access-Control-Allow-Credentials', 'true')
        response.setHeader('Access-Control-Allow-Methods', 'GET,HEAD,OPTIONS')
        response.setHeader('Access-Control-Allow-Headers', request.headers['access-control-request-headers'] || 'Content-Type, X-W7Panel-Token')
        response.setHeader('Access-Control-Allow-Private-Network', 'true')
        if (request.method === 'OPTIONS') return response.end()
        next()
      })
    },
  }, {
    name: 'local-entry-base',
    transformIndexHtml(html) {
      return command === 'serve' ? html.replace('<head>', `<head><base href="${origin}/">`) : html
    },
  }],
  resolve: { alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) } },
  server: {
    host: '0.0.0.0',
    port,
    strictPort: true,
    origin,
    cors: false,
    proxy: { '/panel-api': { target: panelOrigin, changeOrigin: true } },
  },
}))
