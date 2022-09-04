import { defineConfig } from 'vite';
import config from "./vite.config";

export default defineConfig({ ...config, build: {
  outDir: "../cmd/web/assets",
}});
