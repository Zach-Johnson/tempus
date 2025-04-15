import { createRouter, createWebHistory } from "vue-router";

const routes = [
    {
        path: "/",
        name: "home",
        component: () => import("@/views/HomeView.vue"),
        meta: { title: "Home" },
    },
    {
        path: "/categories",
        name: "categories",
        component: () => import("@/views/categories/CategoriesView.vue"),
        meta: { title: "Categories" },
    },
    {
        path: "/categories/:id",
        name: "category-detail",
        component: () => import("@/views/categories/CategoryDetailView.vue"),
        meta: { title: "Category Detail" },
        props: true,
    },
    {
        path: "/tags",
        name: "tags",
        component: () => import("@/views/tags/TagsView.vue"),
        meta: { title: "Tags" },
    },
    {
        path: "/exercises",
        name: "exercises",
        component: () => import("@/views/exercises/ExercisesView.vue"),
        meta: { title: "Exercises" },
    },
    {
        path: "/exercises/:id",
        name: "exercise-detail",
        component: () => import("@/views/exercises/ExerciseDetailView.vue"),
        meta: { title: "Exercise Detail" },
        props: true,
    },
    {
        path: "/sessions",
        name: "sessions",
        component: () => import("@/views/sessions/SessionsView.vue"),
        meta: { title: "Practice Sessions" },
    },
    {
        path: "/sessions/new",
        name: "new-session",
        component: () => import("@/views/sessions/NewSessionView.vue"),
        meta: { title: "New Practice Session" },
    },
    {
        path: "/sessions/:id",
        name: "session-detail",
        component: () => import("@/views/sessions/SessionDetailView.vue"),
        meta: { title: "Session Detail" },
        props: true,
    },
    {
        path: "/:pathMatch(.*)*",
        name: "not-found",
        component: () => import("@/views/NotFoundView.vue"),
        meta: { title: "404 - Not Found" },
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
    scrollBehavior(to, from, savedPosition) {
        if (savedPosition) {
            return savedPosition;
        } else {
            return {
                top: 0,
            };
        }
    },
});

// Update document title based on route
router.beforeEach((to, from, next) => {
    document.title = `${to.meta.title} | Tempus` || "Tempus";
    next();
});

export default router;
