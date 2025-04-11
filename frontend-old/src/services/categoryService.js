import api from "./api.js";

const categoryService = {
    /**
     * Get list of categories with pagination
     * @param {Object} params - { page, page_size }
     * @returns {Promise} Promise with categories data
     */
    getCategories(params = {
        page: 1,
        page_size: 10,
    }) {
        return api.get("/v1/categories", { params });
    },

    /**
     * Get a category by ID
     * @param {Number} id - Category ID
     * @returns {Promise} Promise with category data
     */
    getCategory(id) {
        return api.get(`/v1/categories/${id}`);
    },

    /**
     * Create a new category
     * @param {Object} category - Category data
     * @returns {Promise} Promise with created category
     */
    createCategory(category) {
        return api.post("/v1/categories", category);
    },

    /**
     * Update an existing category
     * @param {Object} category - Category data with ID
     * @returns {Promise} Promise with updated category
     */
    updateCategory(category) {
        return api.put(`/v1/categories/${category.id}`, category);
    },

    /**
     * Delete a category
     * @param {Number} id - Category ID
     * @returns {Promise} Promise with operation result
     */
    deleteCategory(id) {
        return api.delete(`/v1/categories/${id}`);
    },
};

export default categoryService;
