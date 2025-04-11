import { format } from "date-fns";

import api from "./api.js";

const historyService = {
    /**
     * Get exercise history with pagination and filters
     * @param {Object} params - { exercise_id, page, page_size, start_date,
     *     end_date }
     * @returns {Promise} Promise with exercise history data
     */
    getExerciseHistory(params = {
        page: 1,
        page_size: 10,
    }) {
        // Format dates for API
        const formattedParams = { ...params };
        if (params.start_date) {
            formattedParams.start_date = format(
                new Date(params.start_date),
                "yyyy-MM-dd'T'HH:mm:ss'Z'",
            );
        }
        if (params.end_date) {
            formattedParams.end_date = format(
                new Date(params.end_date),
                "yyyy-MM-dd'T'HH:mm:ss'Z'",
            );
        }

        return api.get("/v1/history", { params: formattedParams });
    },

    /**
     * Create a new exercise history entry
     * @param {Object} historyEntry - Exercise history data
     * @returns {Promise} Promise with created history entry
     */
    createExerciseHistory(historyEntry) {
        // Format date for API
        const formattedEntry = { ...historyEntry };
        if (historyEntry.date) {
            formattedEntry.date = format(
                new Date(historyEntry.date),
                "yyyy-MM-dd'T'HH:mm:ss'Z'",
            );
        }

        return api.post("/v1/history", formattedEntry);
    },

    /**
     * Delete an exercise history entry
     * @param {Number} id - History entry ID
     * @returns {Promise} Promise with operation result
     */
    deleteExerciseHistory(id) {
        return api.delete(`/v1/history/${id}`);
    },
};

export default historyService;
