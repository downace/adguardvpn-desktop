import * as path from "node:path";
import { quasar, transformAssetUrls } from "@quasar/vite-plugin";
import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue({
      template: { transformAssetUrls }
    }),

    quasar()
  ],
  resolve: {
    alias: {
      "@/runtime": path.resolve(__dirname, "./wailsjs/runtime"),
      "@/go": path.resolve(__dirname, "./wailsjs/go"),
      "@": path.resolve(__dirname, "./src"),
    },
  },
})
