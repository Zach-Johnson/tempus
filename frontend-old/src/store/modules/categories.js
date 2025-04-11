import categoryService from "@/services/categoryService";
import { defineStore } from "pinia";

export const useCategoriesStore = defineStore("categories", {
    state: () => ({
        categories: [],
        currentCategory: null,
        totalCount: 0,
        loading: false,
        error: null,
    }),

    getters: {
        getCategoryById: (state) => (id) => {
            return state.categories.find((category) => category.id === id);
        },

        categoriesOptions: (state) => {
            return state.categories.map(
                (category) => ({ value: category.id, label: category.name }),
            );
        },
    },

    actions: {
        async fetchCategories(params = {
            page: 1,
            page_size: 10,
        }) {
            this.loading = true;
            this.error = null;

            try {
                const response = await categoryService.getCategories(params);
                this.categories = response.categories || [];
                this.totalCount = response.total_count || 0;
            } catch (error) {
                this.error = error.message || "Failed to fetch categories";
                console.error("Error fetching categories:", error);
            } finally {
                this.loading = false;
            }
        },

        async fetchCategoryById(id) {
            this.loading = true;
            this.error = null;

            try {
                const category = await categoryService.getCategory(id);
                this.currentCategory = category;
                return category;
            } catch (error) {
                this.error = error.message || `Failed to fetch category ${id}`;
                console.error(`Error fetching category ${id}:`, error);
                return null;
            } finally {
                this.loading = false;
            }
        },

        async createCategory(categoryData) {
            this.loading = true;
            this.error = null;

            try {
                const category = await categoryService.createCategory(
                    categoryData,
                );
                this.categories.push(category);
                this.totalCount++;
                return category;
            } catch (error) {
                this.error = error.message || "Failed to create category";
                console.error("Error creating category:", error);
                return null;
            } finally {
                this.loading = false;
            }
        },

        async updateCategory(categoryData) {
            this.loading = true;
            this.error = null;

            try {
                const updatedCategory = await categoryService.updateCategory(
                    categoryData,
                );

                // Update in the list
                const index = this.categories.findIndex((c) =>
                    c.id === updatedCategory.id
                );
                if (index !== -1) {
                    this.categories[index] = updatedCategory;
                }

                // Update current category if it's the same one
                if (
                    this.currentCategory &&
                    this.currentCategory.id === updatedCategory.id
                ) {
                    this.currentCategory = updatedCategory;
                }

                return updatedCategory;
            } catch (error) {
                this.error = error.message || "Failed to update category";
                console.error("Error updating category:", error);
                return null;
            } finally {
                this.loading = false;
            }
        },

        async deleteCategory(id) {
            this.loading = true;
            this.error = null;

            try {
                await categoryService.deleteCategory(id);

                // Remove from the list
                this.categories = this.categories.filter((c) => c.id !== id);
                this.totalCount--;

                // Clear current category if it's the same one
                if (this.currentCategory && this.currentCategory.id === id) {
                    this.currentCategory = null;
                }

                return true;
            } catch (error) {
                this.error = error.message || `Failed to delete category ${id}`;
                console.error(`Error deleting category ${id}:`, error);
                return false;
            } finally {
                this.loading = false;
            }
        },
    },
});
