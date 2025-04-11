// Vuetify
import "vuetify/styles";
// Material Design Icons
import "@mdi/font/css/materialdesignicons.css";
// Global styles
import "./assets/main.scss";

import { createPinia } from "pinia";
import { createApp } from "vue";
import { createVuetify } from "vuetify";
import * as components from "vuetify/components";
import * as directives from "vuetify/directives";

import App from "./App.vue";
import router from "./router/index.js";

const vuetify = createVuetify({
    components,
    directives,
    icons: {
        defaultSet: "mdi",
    },
    theme: {
        defaultTheme: "light",
        themes: {
            light: {
                dark: false,
                colors: {
                    primary: "#1976D2",
                    secondary: "#424242",
                    accent: "#82B1FF",
                    error: "#FF5252",
                    info: "#2196F3",
                    success: "#4CAF50",
                    warning: "#FFC107",
                    background: "#F5F5F5",
                },
            },
            dark: {
                dark: true,
                colors: {
                    primary: "#2196F3",
                    secondary: "#616161",
                    accent: "#82B1FF",
                    error: "#FF5252",
                    info: "#2196F3",
                    success: "#4CAF50",
                    warning: "#FFC107",
                    background: "#121212",
                },
            },
        },
    },
});

// Create and mount the app
const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(vuetify);

app.mount("#app");
