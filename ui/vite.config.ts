import { defineConfig } from 'vite';
import preact from '@preact/preset-vite';

export default defineConfig({
  plugins: [preact()],
  build: {
    outDir: '../content',
    emptyOutDir: true,
  },
  server: {
    proxy: {
      '/ws': 'http://localhost:8080',
    },
  },
});
