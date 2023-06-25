import { defineConfig, UserConfigExport } from 'vite';
import vue from '@vitejs/plugin-vue';

// https://vitejs.dev/config/
export const config: UserConfigExport = {
  plugins: [vue()],
  server: {
    proxy: {
      '/api': 'https://pediaroute.com',
    },
  },
};
export default defineConfig(config);
