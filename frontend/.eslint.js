module.exports = {
    root: true,
    env: {
        node: true,
    },
    extends: [
        "plugin:vue/vue3-essential",
        "eslint:recommended",
    ],
    parserOptions: {
        parser: "babel-eslint",
    },
    rules: {
        "import/no-unresolved": [2, { ignore: ["^@/"] }], // This ignores paths starting with @/
        "no-console": process.env.NODE_ENV === "prod" ? "warn" : "off",
        "no-debugger": process.env.NODE_ENV === "prod" ? "warn" : "off",
    },
};
