import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // Load env file based on `mode` in the current working directory.
  // Set the third parameter to '' to load all env regardless of the `VITE_` prefix.
  const env = loadEnv(mode, process.cwd(), '')

  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    server: {
      proxy: {
        // Proxy API requests to the backend
        '/api': {
          target: env.VITE_PROXY_TARGET || 'http://localhost:8080',
          changeOrigin: true,
          secure: false,
          // Rewrite path if needed: /api/login -> /api/v1/login (if backend doesn't use /api prefix)
          // However, if backend already uses /api, we don't need rewrite.
          // In this project, backend seems to use /api/v1 or just /
          // Let's assume the backend has /api/v1 prefix as seen in previous conversations
          // rewrite: (path) => path.replace(/^\/api/, ''), 
        },
        // Proxy WebSocket requests
        '/ws': {
          target: env.VITE_PROXY_TARGET || 'http://localhost:8080',
          ws: true,
          changeOrigin: true,
          secure: false,
          rewrite: (path) => path.replace(/^\/ws/, '/api/v1/ws'),
        },
      },
    },
  }
})
