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
                    return new Date(b.start_time) - new Date(a.start_time);
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
                            (ex) => ex.exercise_id === exerciseId,
                        ),
                );
        },
    );

    const totalPracticeTime = computed(() => {
        return sessions.value.reduce((total, session) => {
            if (session.start_time && session.end_time) {
                const startTime = new Date(session.start_time);
                const endTime = new Date(session.end_time);
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
        } catch (err) {
            error.value = err.message || "Failed to fetch practice sessions";
            console.error("Error fetching practice sessions:", err);
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
        } catch (err) {
            error.value = err.message ||
                `Failed to fetch practice session with ID ${id}`;
            console.error(`Error fetching practice session ${id}:`, err);
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
            sessions.value.push(newSession);
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
                `Failed to update practice session with ID ${id}`;
            console.error(`Error updating practice session ${id}:`, err);
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
                `Failed to delete practice session with ID ${id}`;
            console.error(`Error deleting practice session ${id}:`, err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function addSessionExercise(sessionId, exerciseData) {
        loading.value = true;
        error.value = null;

        try {
            const response = await sessionsAPI.addExercise(
                sessionId,
                exerciseData,
            );

            // Update session with new exercise if it's in the sessions array
            const session = sessions.value.find((s) => s.id === sessionId);
            if (session) {
                if (!session.exercises) {
                    session.exercises = [];
                }
                session.exercises.push(response.data);
            }

            // Update currentSession if it's the relevant session
            if (currentSession.value && currentSession.value.id === sessionId) {
                if (!currentSession.value.exercises) {
                    currentSession.value.exercises = [];
                }
                currentSession.value.exercises.push(response.data);
            }

            return response.data;
        } catch (err) {
            error.value = err.message ||
                `Failed to add exercise to session with ID ${sessionId}`;
            console.error(
                `Error adding exercise to session ${sessionId}:`,
                err,
            );
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function updateSessionExercise(
        sessionId,
        exerciseId,
        exerciseData,
        updateMask,
    ) {
        loading.value = true;
        error.value = null;

        try {
            const response = await sessionsAPI.updateExercise(
                sessionId,
                exerciseId,
                exerciseData,
                updateMask,
            );
            const updatedExercise = response.data;

            // Update exercise in session if it's in the sessions array
            const session = sessions.value.find((s) => s.id === sessionId);
            if (session && session.exercises) {
                const index = session.exercises.findIndex((e) =>
                    e.id === exerciseId
                );
                if (index !== -1) {
                    session.exercises[index] = updatedExercise;
                }
            }

            // Update exercise in currentSession if applicable
            if (
                currentSession.value && currentSession.value.id === sessionId &&
                currentSession.value.exercises
            ) {
                const index = currentSession.value.exercises.findIndex(
                    (e) => e.id === exerciseId,
                );
                if (index !== -1) {
                    currentSession.value.exercises[index] = updatedExercise;
                }
            }

            return updatedExercise;
        } catch (err) {
            error.value = err.message ||
                `Failed to update exercise ${exerciseId} in session ${sessionId}`;
            console.error(
                `Error updating exercise ${exerciseId} in session ${sessionId}:`,
                err,
            );
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function deleteSessionExercise(exerciseId) {
        loading.value = true;
        error.value = null;

        try {
            await sessionsAPI.deleteExercise(exerciseId);

            // Remove exercise from sessions in the array
            sessions.value.forEach((session) => {
                if (session.exercises) {
                    session.exercises = session.exercises.filter((ex) =>
                        ex.id !== exerciseId
                    );
                }
            });

            // Remove exercise from currentSession if applicable
            if (currentSession.value && currentSession.value.exercises) {
                currentSession.value.exercises = currentSession.value.exercises
                    .filter(
                        (ex) => ex.id !== exerciseId,
                    );
            }
        } catch (err) {
            error.value = err.message ||
                `Failed to delete exercise with ID ${exerciseId} from session`;
            console.error(
                `Error deleting exercise ${exerciseId} from session:`,
                err,
            );
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
        addSessionExercise,
        updateSessionExercise,
        deleteSessionExercise,
        fetchPracticeStats,
    };
});
