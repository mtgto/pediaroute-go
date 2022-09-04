
import { defineConfig, UserConfigExport } from 'vite';
import vue from '@vitejs/plugin-vue';
import vueI18n from '@intlify/vite-plugin-vue-i18n';

// https://vitejs.dev/config/
export const config: UserConfigExport = {
  plugins: [vue(), vueI18n()],
  server: {
    proxy: {
      '/api': 'https://pediaroute.com',
    }
  },
};
export default defineConfig(config);
