import { categoriesAPI } from "@/services/api.js";
import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useCategoriesStore = defineStore("categories", () => {
    // State
    const categories = ref([]);
    const loading = ref(false);
    const error = ref(null);
    const currentCategory = ref(null);
    const currentCategoryLoading = ref(false);

    // Getters
    const categoriesSorted = computed(
        () => {
            return [...categories.value].sort(
                (a, b) => a.name.localeCompare(b.name),
            );
        },
    );

    const categoryById = computed(() => {
        return (id) => categories.value.find((cat) => cat.id === id);
    });

    // Actions
    async function fetchCategories() {
        loading.value = true;
        error.value = null;

        try {
            const response = await categoriesAPI.getAll();
            categories.value = response.data.categories || [];
        } catch (err) {
            error.value = err.message || "Failed to fetch categories";
            console.error("Error fetching categories:", err);
        } finally {
            loading.value = false;
        }
    }

    async function fetchCategory(id) {
        currentCategoryLoading.value = true;
        error.value = null;

        try {
            const response = await categoriesAPI.get(id);
            currentCategory.value = response.data;

            // Also update the category in the categories array if it exists
            const index = categories.value.findIndex((c) => c.id === id);
            if (index !== -1) {
                categories.value[index] = response.data;
            } else {
                categories.value.push(response.data);
            }
        } catch (err) {
            error.value = err.message ||
                `Failed to fetch category with ID ${id}`;
            console.error(`Error fetching category ${id}:`, err);
        } finally {
            currentCategoryLoading.value = false;
        }
    }

    async function createCategory(categoryData) {
        loading.value = true;
        error.value = null;

        try {
            const response = await categoriesAPI.create(categoryData);
            const newCategory = response.data;
            categories.value.push(newCategory);
            return newCategory;
        } catch (err) {
            error.value = err.message || "Failed to create category";
            console.error("Error creating category:", err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function updateCategory(id, categoryData, updateMask) {
        loading.value = true;
        error.value = null;

        try {
            const response = await categoriesAPI.update(
                id,
                categoryData,
                updateMask,
            );
            const updatedCategory = response.data;

            // Update in the categories array
            const index = categories.value.findIndex((c) => c.id === id);
            if (index !== -1) {
                categories.value[index] = updatedCategory;
            }

            // Update currentCategory if it's the one being edited
            if (currentCategory.value && currentCategory.value.id === id) {
                currentCategory.value = updatedCategory;
            }

            return updatedCategory;
        } catch (err) {
            error.value = err.message ||
                `Failed to update category with ID ${id}`;
            console.error(`Error updating category ${id}:`, err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function deleteCategory(id) {
        loading.value = true;
        error.value = null;

        try {
            await categoriesAPI.delete(id);

            // Remove from the categories array
            categories.value = categories.value.filter((c) => c.id !== id);

            // Clear currentCategory if it's the one being deleted
            if (currentCategory.value && currentCategory.value.id === id) {
                currentCategory.value = null;
            }
        } catch (err) {
            error.value = err.message ||
                `Failed to delete category with ID ${id}`;
            console.error(`Error deleting category ${id}:`, err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    return {
        // State
        categories,
        loading,
        error,
        currentCategory,
        currentCategoryLoading,

        // Getters
        categoriesSorted,
        categoryById,

        // Actions
        fetchCategories,
        fetchCategory,
        createCategory,
        updateCategory,
        deleteCategory,
    };
});
