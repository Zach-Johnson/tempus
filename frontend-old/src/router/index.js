import CategoriesView from "@/views/CategoriesView.vue";
// Import views
import Dashboard from "@/views/Dashboard.vue";
import ExercisesView from "@/views/ExercisesView.vue";
// import HistoryView from "@/views/HistoryView.vue";
import PracticeSessionsView from "@/views/PracticeSessionsView.vue";
import TagsView from "@/views/TagsView.vue";
import { createRouter, createWebHistory } from "vue-router";

const routes = [
    { path: "/", name: "Dashboard", component: Dashboard },
    { path: "/categories", name: "Categories", component: CategoriesView },
    { path: "/tags", name: "Tags", component: TagsView },
    {
        path: "/exercises",
        name: "Exercises",
        component: ExercisesView,
        children: [{
            path: ":id",
            name: "ExerciseDetails",
            component: () =>
                import("@/components/exercises/ExerciseDetails.vue"),
            props: (route) => ({ exerciseId: Number(route.params.id) }),
        }],
    },
    {
        path: "/practice",
        name: "PracticeSessions",
        component: PracticeSessionsView,
        children: [{
            path: ":id",
            name: "SessionDetails",
            component: () => import("@/components/practice/SessionDetails.vue"),
            props: (route) => ({ sessionId: Number(route.params.id) }),
        }],
    },
    // { path: "/history", name: "History", component: HistoryView },
    { path: "/:pathMatch(.*)*", redirect: "/" },
];

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes,
});

export default router;
