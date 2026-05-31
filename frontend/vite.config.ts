import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { fileURLToPath, URL } from 'node:url'
// @ts-ignore
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [react(), wails('./bindings')],
  resolve: {
    dedupe: ['react', 'react-dom'],
    alias: {
      // @ts-ignore
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      // @ts-ignore
      '@bindings': fileURLToPath(new URL('./bindings', import.meta.url)),
    },
  },
  // @ts-ignore
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/__tests__/setup.ts'],
  },
})
