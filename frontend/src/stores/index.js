// Import all store files
import { useAppStore } from "./app.js";
import { useCategoriesStore } from "./categories.js";
import { useExercisesStore } from "./exercises.js";
import { useHistoryStore } from "./history.js";
import { useSessionsStore } from "./sessions.js";
import { useTagsStore } from "./tags.js";

// Export a function to get all stores for easy access
export function useStores() {
    const appStore = useAppStore();
    const categoriesStore = useCategoriesStore();
    const tagsStore = useTagsStore();
    const exercisesStore = useExercisesStore();
    const sessionsStore = useSessionsStore();
    const historyStore = useHistoryStore();

    return {
        app: appStore,
        categories: categoriesStore,
        tags: tagsStore,
        exercises: exercisesStore,
        sessions: sessionsStore,
        history: historyStore,
    };
}

// Individual exports for when only specific stores are needed
export {
    useAppStore,
    useCategoriesStore,
    useExercisesStore,
    useHistoryStore,
    useSessionsStore,
    useTagsStore,
};
