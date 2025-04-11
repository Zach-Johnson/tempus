import api from "./api.js";

const exerciseService = {
    /**
     * Get list of exercises with pagination and filters
     * @param {Object} params - { page, page_size, category_ids, tag_ids,
     *     search_term }
     * @returns {Promise} Promise with exercises data
     */
    getExercises(params = {
        page: 1,
        page_size: 10,
    }) {
        return api.get("/v1/exercises", { params });
    },

    /**
     * Get an exercise by ID
     * @param {Number} id - Exercise ID
     * @returns {Promise} Promise with exercise data
     */
    getExercise(id) {
        return api.get(`/v1/exercises/${id}`);
    },

    /**
     * Create a new exercise
     * @param {Object} exercise - Exercise data
     * @returns {Promise} Promise with created exercise
     */
    createExercise(exercise) {
        return api.post("/v1/exercises", exercise);
    },

    /**
     * Update an existing exercise
     * @param {Object} exercise - Exercise data with ID
     * @returns {Promise} Promise with updated exercise
     */
    updateExercise(exercise) {
        return api.put(`/v1/exercises/${exercise.id}`, exercise);
    },

    /**
     * Delete an exercise
     * @param {Number} id - Exercise ID
     * @returns {Promise} Promise with operation result
     */
    deleteExercise(id) {
        return api.delete(`/v1/exercises/${id}`);
    },

    /**
     * Add a link to an exercise
     * @param {Object} link - Link data with exercise_id
     * @returns {Promise} Promise with created link
     */
    addExerciseLink(link) {
        return api.post("/v1/exercise-links", link);
    },

    /**
     * Delete an exercise link
     * @param {Number} id - Link ID
     * @returns {Promise} Promise with operation result
     */
    deleteExerciseLink(id) {
        return api.delete(`/v1/exercise-links/${id}`);
    },

    /**
     * Add an exercise to a category
     * @param {Number} exerciseId - Exercise ID
     * @param {Number} categoryId - Category ID
     * @returns {Promise} Promise with operation result
     */
    addExerciseToCategory(exerciseId, categoryId) {
        return api.post(
            "/v1/exercise-categories",
            { exercise_id: exerciseId, category_id: categoryId },
        );
    },

    /**
     * Remove an exercise from a category
     * @param {Number} exerciseId - Exercise ID
     * @param {Number} categoryId - Category ID
     * @returns {Promise} Promise with operation result
     */
    removeExerciseFromCategory(exerciseId, categoryId) {
        return api.delete(
            `/v1/exercise-categories`,
            { data: { exercise_id: exerciseId, category_id: categoryId } },
        );
    },

    /**
     * Add a tag to an exercise
     * @param {Number} exerciseId - Exercise ID
     * @param {Number} tagId - Tag ID
     * @returns {Promise} Promise with operation result
     */
    addTagToExercise(exerciseId, tagId) {
        return api.post(`/v1/exercises/${exerciseId}/tags/${tagId}`);
    },

    /**
     * Remove a tag from an exercise
     * @param {Number} exerciseId - Exercise ID
     * @param {Number} tagId - Tag ID
     * @returns {Promise} Promise with operation result
     */
    removeTagFromExercise(exerciseId, tagId) {
        return api.delete(`/v1/exercises/${exerciseId}/tags/${tagId}`);
    },
};

export default exerciseService;
