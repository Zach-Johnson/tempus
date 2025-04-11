import vue from "@vitejs/plugin-vue";
import path from "path";
import { defineConfig } from "vite";

export default defineConfig({
    plugins: [vue()],
    server: {
        proxy: {
            "/api": {
                target: process.env.VITE_API_URL || "http://localhost:8080",
                changeOrigin: true,
                rewrite: (path) => path.replace(/^\/api/, ""),
            },
        },
    },
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "./src"),
        },
    },
});
