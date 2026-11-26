import { resolve } from "node:path";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

const BACKEND_ORIGIN = process.env.OMEG4MON_BACKEND_URL ?? "http://127.0.0.1:20128";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": resolve(import.meta.dirname, "src"),
    },
  },
  build: {
    outDir: "out",
    emptyOutDir: true,
    sourcemap: false,
    rolldownOptions: {
      output: {
        codeSplitting: {
          groups: [
            {
              name: "vendor-react",
              test: /node_modules[\\/](react|react-dom|react-router|scheduler)[\\/]/,
              priority: 20,
            },
            {
              name: "vendor-query",
              test: /node_modules[\\/](@tanstack[\\/]react-query|openapi-fetch|zustand)[\\/]/,
              priority: 15,
            },
          ],
        },
      },
    },
  },
  server: {
    host: "127.0.0.1",
    port: 20127,
    strictPort: true,
    proxy: {
      "/api": {
        target: BACKEND_ORIGIN,
        changeOrigin: true,
      },
      "/v1": {
        target: BACKEND_ORIGIN,
        changeOrigin: true,
      },
      "/v1beta": {
        target: BACKEND_ORIGIN,
        changeOrigin: true,
      },
    },
  },
  preview: {
    host: "127.0.0.1",
    port: 20127,
    strictPort: true,
  },
});
