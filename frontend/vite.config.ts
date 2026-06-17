import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    // Proxy /api/* to the Go backend during development so the browser
    // never makes cross-origin requests (avoids CORS issues).
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        // No rewrite — Go routes now include the /api prefix.
      },
    },
  },
})
