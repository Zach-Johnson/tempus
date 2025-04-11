import api from "./api.js";

const tagService = {
    /**
     * Get list of tags with pagination
     * @param {Object} params - { page, page_size }
     * @returns {Promise} Promise with tags data
     */
    getTags(params = {
        page: 1,
        page_size: 10,
    }) {
        return api.get("/v1/tags", { params });
    },

    /**
     * Get a tag by ID
     * @param {Number} id - Tag ID
     * @returns {Promise} Promise with tag data
     */
    getTag(id) {
        return api.get(`/v1/tags/${id}`);
    },

    /**
     * Create a new tag
     * @param {Object} tag - Tag data
     * @returns {Promise} Promise with created tag
     */
    createTag(tag) {
        return api.post("/v1/tags", tag);
    },

    /**
     * Delete a tag
     * @param {Number} id - Tag ID
     * @returns {Promise} Promise with operation result
     */
    deleteTag(id) {
        return api.delete(`/v1/tags/${id}`);
    },
};

export default tagService;
