import { format } from "date-fns";

import api from "./api.js";

const sessionService = {
    /**
     * Get list of practice sessions with pagination and date filters
     * @param {Object} params - { page, page_size, start_date, end_date }
     * @returns {Promise} Promise with practice sessions data
     */
    getPracticeSessions(params = {
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

        return api.get("/v1/sessions", { params: formattedParams });
    },

    /**
     * Get a practice session by ID
     * @param {Number} id - Practice session ID
     * @returns {Promise} Promise with practice session data
     */
    getPracticeSession(id) {
        return api.get(`/v1/sessions/${id}`);
    },

    /**
     * Create a new practice session
     * @param {Object} session - Practice session data
     * @returns {Promise} Promise with created practice session
     */
    createPracticeSession(session) {
        // Format date for API
        const formattedSession = { ...session };
        if (session.date) {
            formattedSession.date = format(
                new Date(session.date),
                "yyyy-MM-dd'T'HH:mm:ss'Z'",
            );
        }

        return api.post("/v1/sessions", formattedSession);
    },

    /**
     * Update an existing practice session
     * @param {Object} session - Practice session data with ID
     * @returns {Promise} Promise with updated practice session
     */
    updatePracticeSession(session) {
        // Format date for API
        const formattedSession = { ...session };
        if (session.date) {
            formattedSession.date = format(
                new Date(session.date),
                "yyyy-MM-dd'T'HH:mm:ss'Z'",
            );
        }

        return api.put(`/v1/sessions/${session.id}`, formattedSession);
    },

    /**
     * Delete a practice session
     * @param {Number} id - Practice session ID
     * @returns {Promise} Promise with operation result
     */
    deletePracticeSession(id) {
        return api.delete(`/v1/sessions/${id}`);
    },

    /**
     * Add an exercise to a practice session
     * @param {Object} sessionExercise - Session exercise data
     * @returns {Promise} Promise with created session exercise
     */
    addExerciseToSession(sessionId, sessionExercise) {
        console.log(sessionExercise);
        return api.post(`/v1/sessions/${sessionId}/exercises`, sessionExercise);
    },

    /**
     * Update a session exercise
     * @param {Object} sessionExercise - Session exercise data with ID
     * @returns {Promise} Promise with updated session exercise
     */
    updateSessionExercise(sessionExercise) {
        return api.put(
            `/v1/sessions/exercises/${sessionExercise.id}`,
            sessionExercise,
        );
    },

    /**
     * Remove an exercise from a session
     * @param {Number} sessionId - Session ID
     * @param {Number} exerciseId - Exercise ID
     * @returns {Promise} Promise with operation result
     */
    removeExerciseFromSession(sessionId, exerciseId) {
        return api.delete(
            `/v1/sessions/${sessionId}/exercises/${exerciseId}`,
        );
    },
};

export default sessionService;
