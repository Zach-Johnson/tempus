// Import Bootstrap CSS and JS
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap/dist/js/bootstrap.bundle.min.js";
// Import global CSS
import "./assets/styles/main.css";

import { createPinia } from "pinia";
import { createApp } from "vue";

import App from "./App.vue";
import router from "./router/index.js";

const app = createApp(App);

// Use Pinia for state management
const pinia = createPinia();
app.use(pinia);

// Use Vue Router
app.use(router);

app.mount("#app");
