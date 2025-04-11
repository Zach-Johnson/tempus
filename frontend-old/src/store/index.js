// Import the createPinia function
import { createPinia } from "pinia";

// Import store modules
import { useCategoriesStore } from "./modules/categories.js";
import { useExercisesStore } from "./modules/exercises.js";
import { useHistoryStore } from "./modules/history.js";
import { useSessionsStore } from "./modules/sessions.js";
import { useTagsStore } from "./modules/tags.js";

// Create and export the Pinia instance
const pinia = createPinia();

export {
    pinia,
    useCategoriesStore,
    useExercisesStore,
    useHistoryStore,
    useSessionsStore,
    useTagsStore,
};
