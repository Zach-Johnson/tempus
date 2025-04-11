import sessionService from "@/services/sessionService";
import { defineStore } from "pinia";

export const useSessionsStore = defineStore("sessions", {
    state: () => ({
        sessions: [],
        currentSession: null,
        totalCount: 0,
        loading: false,
        error: null,
        dateRange: { startDate: null, endDate: null },
    }),

    getters: {
        getSessionById: (state) => (id) => {
            return state.sessions.find((session) => session.id === id);
        },

        totalDuration: (state) => {
            return state.sessions.reduce(
                (total, session) => total + (session.duration_minutes || 0),
                0,
            );
        },

        sessionsByDate: (state) => {
            const sessionMap = {};

            state.sessions.forEach((session) => {
                const date = new Date(session.date).toISOString().split("T")[0];
                if (!sessionMap[date]) {
                    sessionMap[date] = [];
                }
                sessionMap[date].push(session);
            });

            return sessionMap;
        },
    },

    actions: {
        setDateRange(startDate, endDate) {
            this.dateRange = { startDate, endDate };
        },

        clearDateRange() {
            this.dateRange = { startDate: null, endDate: null };
        },

        async fetchSessions(params = {
            page: 1,
            page_size: 10,
        }) {
            this.loading = true;
            this.error = null;

            try {
                // Apply date filters
                const queryParams = {
                    ...params,
                    start_date: this.dateRange.startDate,
                    end_date: this.dateRange.endDate,
                };

                const response = await sessionService.getPracticeSessions(
                    queryParams,
                );
                this.sessions = response.sessions || [];
                this.totalCount = response.total_count || 0;
            } catch (error) {
                this.error = error.message ||
                    "Failed to fetch practice sessions";
                console.error("Error fetching practice sessions:", error);
            } finally {
                this.loading = false;
            }
        },

        async fetchSessionById(id) {
            this.loading = true;
            this.error = null;

            try {
                const session = await sessionService.getPracticeSession(id);
                this.currentSession = session;
                return session;
            } catch (error) {
                this.error = error.message ||
                    `Failed to fetch practice session ${id}`;
                console.error(`Error fetching practice session ${id}:`, error);
                return null;
            } finally {
                this.loading = false;
            }
        },

        async createSession(sessionData) {
            this.loading = true;
            this.error = null;

            try {
                const session = await sessionService.createPracticeSession(
                    sessionData,
                );
                this.sessions.push(session);
                this.totalCount++;
                return session;
            } catch (error) {
                this.error = error.message ||
                    "Failed to create practice session";
                console.error("Error creating practice session:", error);
                return null;
            } finally {
                this.loading = false;
            }
        },

        async updateSession(sessionData) {
            this.loading = true;
            this.error = null;

            try {
                const updatedSession = await sessionService
                    .updatePracticeSession(sessionData);

                // Update in the list
                const index = this.sessions.findIndex((s) =>
                    s.id === updatedSession.id
                );
                if (index !== -1) {
                    this.sessions[index] = updatedSession;
                }

                // Update current session if it's the same one
                if (
                    this.currentSession &&
                    this.currentSession.id === updatedSession.id
                ) {
                    this.currentSession = updatedSession;
                }

                return updatedSession;
            } catch (error) {
                this.error = error.message ||
                    "Failed to update practice session";
                console.error("Error updating practice session:", error);
                return null;
            } finally {
                this.loading = false;
            }
        },

        async deleteSession(id) {
            this.loading = true;
            this.error = null;

            try {
                await sessionService.deletePracticeSession(id);

                // Remove from the list
                this.sessions = this.sessions.filter((s) => s.id !== id);
                this.totalCount--;

                // Clear current session if it's the same one
                if (this.currentSession && this.currentSession.id === id) {
                    this.currentSession = null;
                }

                return true;
            } catch (error) {
                this.error = error.message ||
                    `Failed to delete practice session ${id}`;
                console.error(`Error deleting practice session ${id}:`, error);
                return false;
            } finally {
                this.loading = false;
            }
        },

        async addExerciseToSession(sessionId, sessionExercise) {
            this.loading = true;
            this.error = null;

            try {
                const newSessionExercise = await sessionService
                    .addExerciseToSession(
                        sessionId,
                        sessionExercise,
                    );

                // Update the session exercises if it's the current session
                if (
                    this.currentSession &&
                    this.currentSession.id === sessionExercise.session_id
                ) {
                    if (!this.currentSession.exercises) {
                        this.currentSession.exercises = [];
                    }
                    this.currentSession.exercises.push(newSessionExercise);

                    // Recalculate total duration if needed
                    if (this.currentSession.duration_minutes) {
                        this.currentSession.duration_minutes +=
                            newSessionExercise.duration_minutes || 0;
                    }
                }

                return newSessionExercise;
            } catch (error) {
                this.error = error.message ||
                    "Failed to add exercise to session";
                console.error("Error adding exercise to session:", error);
                return null;
            } finally {
                this.loading = false;
            }
        },

        async updateSessionExercise(sessionExercise) {
            this.loading = true;
            this.error = null;

            try {
                const updatedSessionExercise = await sessionService
                    .updateSessionExercise(sessionExercise);

                // Update in the current session if applicable
                if (this.currentSession && this.currentSession.exercises) {
                    const index = this.currentSession.exercises.findIndex(
                        (se) => se.id === updatedSessionExercise.id,
                    );

                    if (index !== -1) {
                        // Track duration difference for session total update
                        const oldDuration = this.currentSession.exercises[index]
                            .duration_minutes ||
                            0;
                        const newDuration =
                            updatedSessionExercise.duration_minutes || 0;
                        const durationDiff = newDuration - oldDuration;

                        // Update the exercise
                        this.currentSession.exercises[index] =
                            updatedSessionExercise;

                        // Update total session duration if needed
                        if (
                            this.currentSession.duration_minutes &&
                            durationDiff !== 0
                        ) {
                            this.currentSession.duration_minutes +=
                                durationDiff;
                        }
                    }
                }

                return updatedSessionExercise;
            } catch (error) {
                this.error = error.message ||
                    "Failed to update session exercise";
                console.error("Error updating session exercise:", error);
                return null;
            } finally {
                this.loading = false;
            }
        },

        async removeExerciseFromSession(sessionId, exerciseId) {
            this.loading = true;
            this.error = null;

            try {
                await sessionService.removeExerciseFromSession(
                    sessionId,
                    exerciseId,
                );

                // Update the current session if applicable
                if (
                    this.currentSession &&
                    this.currentSession.id === sessionId
                ) {
                    // Find the exercise to remove
                    const exerciseToRemove = this.currentSession.exercises.find(
                        (se) => se.exercise_id === exerciseId,
                    );

                    // Subtract its duration from the session total if needed
                    if (
                        exerciseToRemove &&
                        this.currentSession.duration_minutes
                    ) {
                        this.currentSession.duration_minutes -=
                            exerciseToRemove.duration_minutes || 0;
                        if (this.currentSession.duration_minutes < 0) {
                            this.currentSession.duration_minutes = 0;
                        }
                    }

                    // Remove the exercise
                    this.currentSession.exercises = this.currentSession
                        .exercises.filter(
                            (se) => se.exercise_id !== exerciseId,
                        );
                }

                return true;
            } catch (error) {
                this.error = error.message ||
                    "Failed to remove exercise from session";
                console.error("Error removing exercise from session:", error);
                return false;
            } finally {
                this.loading = false;
            }
        },
    },
});
