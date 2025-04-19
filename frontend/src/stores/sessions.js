import { sessionsAPI } from "@/services/api.js";
import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useSessionsStore = defineStore("sessions", () => {
    // State
    const sessions = ref([]);
    const loading = ref(false);
    const error = ref(null);
    const currentSession = ref(null);
    const currentSessionLoading = ref(false);
    const practiceStats = ref(null);
    const statsLoading = ref(false);

    // Getters
    const sessionsSortedByDate = computed(
        () => {
            return [...sessions.value].sort(
                (a, b) => { // Sort by start_time in descending order (newest
                    // first)
                    return new Date(b.startTime) - new Date(a.startTime);
                },
            );
        },
    );

    const sessionById = computed(
        () => {
            return (id) => sessions.value.find((session) => session.id === id);
        },
    );

    const sessionsByExercise = computed(
        () => {
            return (exerciseId) =>
                sessions.value.filter(
                    (session) =>
                        session.exercises &&
                        session.exercises.some(
                            (ex) => ex.exerciseId === exerciseId,
                        ),
                );
        },
    );

    const totalPracticeTime = computed(() => {
        return sessions.value.reduce((total, session) => {
            if (session.startTime && session.endTime) {
                const startTime = new Date(session.startTime);
                const endTime = new Date(session.endTime);
                const durationMinutes = (endTime - startTime) / (1000 * 60);
                return total + durationMinutes;
            }
            return total;
        }, 0);
    });

    // Actions
    async function fetchSessions(params = {}) {
        loading.value = true;
        error.value = null;

        try {
            const response = await sessionsAPI.getAll(params);
            sessions.value = response.data.sessions || [];
            return response.data;
        } catch (err) {
            error.value = err.message || "Failed to fetch practice sessions";
            console.error("Error fetching practice sessions:", err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function fetchSession(id) {
        currentSessionLoading.value = true;
        error.value = null;

        try {
            const response = await sessionsAPI.get(id);
            currentSession.value = response.data;

            // Also update the session in the sessions array if it exists
            const index = sessions.value.findIndex((s) => s.id === id);
            if (index !== -1) {
                sessions.value[index] = response.data;
            } else {
                sessions.value.push(response.data);
            }

            return response.data;
        } catch (err) {
            error.value = err.message ||
                `Failed to fetch practice session with ID ${id}`;
            console.error(`Error fetching practice session ${id}:`, err);
            throw err;
        } finally {
            currentSessionLoading.value = false;
        }
    }

    async function createSession(sessionData) {
        loading.value = true;
        error.value = null;

        try {
            const response = await sessionsAPI.create(sessionData);
            const newSession = response.data;
            sessions.value.unshift(
                newSession,
            ); // Add to the beginning of the array
            return newSession;
        } catch (err) {
            error.value = err.message || "Failed to create practice session";
            console.error("Error creating practice session:", err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function updateSession(id, sessionData, updateMask) {
        loading.value = true;
        error.value = null;

        try {
            const response = await sessionsAPI.update(
                id,
                sessionData,
                updateMask,
            );
            const updatedSession = response.data;

            // Update in the sessions array
            const index = sessions.value.findIndex((s) => s.id === id);
            if (index !== -1) {
                sessions.value[index] = updatedSession;
            }

            // Update currentSession if it's the one being edited
            if (currentSession.value && currentSession.value.id === id) {
                currentSession.value = updatedSession;
            }

            return updatedSession;
        } catch (err) {
            error.value = err.message ||
                `Failed to update session with ID ${id}`;
            console.error(`Error updating session ${id}:`, err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function deleteSession(id) {
        loading.value = true;
        error.value = null;

        try {
            await sessionsAPI.delete(id);

            // Remove from the sessions array
            sessions.value = sessions.value.filter((s) => s.id !== id);

            // Clear currentSession if it's the one being deleted
            if (currentSession.value && currentSession.value.id === id) {
                currentSession.value = null;
            }
        } catch (err) {
            error.value = err.message ||
                `Failed to delete session with ID ${id}`;
            console.error(`Error deleting session ${id}:`, err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function fetchPracticeStats(params = {}) {
        statsLoading.value = true;
        error.value = null;

        try {
            const response = await sessionsAPI.getStats(params);
            practiceStats.value = response.data;
            return response.data;
        } catch (err) {
            error.value = err.message || "Failed to fetch practice statistics";
            console.error("Error fetching practice statistics:", err);
            throw err;
        } finally {
            statsLoading.value = false;
        }
    }

    async function checkForActiveSession() {
        try {
            const response = await sessionsAPI.getAll({ active: true });
            return response.data.sessions && response.data.sessions.length > 0
                ? response.data.sessions[0]
                : null;
        } catch (err) {
            console.error("Error checking for active sessions:", err);
            return null;
        }
    }

    // Return the store
    return {
        // State
        sessions,
        loading,
        error,
        currentSession,
        currentSessionLoading,
        practiceStats,
        statsLoading,

        // Getters
        sessionsSortedByDate,
        sessionById,
        sessionsByExercise,
        totalPracticeTime,

        // Actions
        fetchSessions,
        fetchSession,
        createSession,
        updateSession,
        deleteSession,
        fetchPracticeStats,
        checkForActiveSession,
    };
});
