const { defineConfig } = require("@vue/cli-service");

module.exports = defineConfig({
    transpileDependencies: true,
    lintOnSave: false, // Disable ESLint completely
    devServer: {
        proxy: {
            "/api": {
                target: process.env.VUE_APP_API_URL || "http://localhost:8081",
                changeOrigin: true,
                pathRewrite: { "^/api": "" },
            },
        },
    },
});
