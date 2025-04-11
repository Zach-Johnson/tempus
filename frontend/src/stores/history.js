import { historyAPI } from "@/services/api.js";
import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useHistoryStore = defineStore("history", () => {
    // State
    const historyEntries = ref([]);
    const loading = ref(false);
    const error = ref(null);
    const currentHistoryEntry = ref(null);
    const currentHistoryLoading = ref(false);

    // Getters
    const entriesSortedByDate = computed(
        () => {
            return [...historyEntries.value].sort(
                (a, b) => { // Sort by start_time in descending order (newest
                    // first)
                    return new Date(b.start_time) - new Date(a.start_time);
                },
            );
        },
    );

    const entriesByExercise = computed(
        () => {
            return (exerciseId) =>
                historyEntries.value.filter(
                    (entry) => entry.exercise_id === exerciseId,
                );
        },
    );

    const entriesBySession = computed(
        () => {
            return (sessionId) =>
                historyEntries.value.filter(
                    (entry) => entry.session_id === sessionId,
                );
        },
    );

    const entriesByDateRange = computed(() => {
        return (startDate, endDate) => {
            const start = startDate ? new Date(startDate) : null;
            const end = endDate ? new Date(endDate) : null;

            return historyEntries.value.filter((entry) => {
                const entryDate = new Date(entry.start_time);

                if (start && end) {
                    return entryDate >= start && entryDate <= end;
                } else if (start) {
                    return entryDate >= start;
                } else if (end) {
                    return entryDate <= end;
                }

                return true;
            });
        };
    });

    // Actions
    async function fetchHistoryEntries(params = {}) {
        loading.value = true;
        error.value = null;

        try {
            const response = await historyAPI.getAll(params);
            historyEntries.value = response.data.history_entries || [];
        } catch (err) {
            error.value = err.message || "Failed to fetch history entries";
            console.error("Error fetching history entries:", err);
        } finally {
            loading.value = false;
        }
    }

    async function fetchHistoryEntry(id) {
        currentHistoryLoading.value = true;
        error.value = null;

        try {
            const response = await historyAPI.get(id);
            currentHistoryEntry.value = response.data;

            // Also update the entry in the historyEntries array if it exists
            const index = historyEntries.value.findIndex((entry) =>
                entry.id === id
            );
            if (index !== -1) {
                historyEntries.value[index] = response.data;
            } else {
                historyEntries.value.push(response.data);
            }
        } catch (err) {
            error.value = err.message ||
                `Failed to fetch history entry with ID ${id}`;
            console.error(`Error fetching history entry ${id}:`, err);
        } finally {
            currentHistoryLoading.value = false;
        }
    }

    async function createHistoryEntry(entryData) {
        loading.value = true;
        error.value = null;

        try {
            const response = await historyAPI.create(entryData);
            const newEntry = response.data;
            historyEntries.value.push(newEntry);
            return newEntry;
        } catch (err) {
            error.value = err.message || "Failed to create history entry";
            console.error("Error creating history entry:", err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function updateHistoryEntry(id, entryData, updateMask) {
        loading.value = true;
        error.value = null;

        try {
            const response = await historyAPI.update(id, entryData, updateMask);
            const updatedEntry = response.data;

            // Update in the historyEntries array
            const index = historyEntries.value.findIndex((entry) =>
                entry.id === id
            );
            if (index !== -1) {
                historyEntries.value[index] = updatedEntry;
            }

            // Update currentHistoryEntry if it's the one being edited
            if (
                currentHistoryEntry.value &&
                currentHistoryEntry.value.id === id
            ) {
                currentHistoryEntry.value = updatedEntry;
            }

            return updatedEntry;
        } catch (err) {
            error.value = err.message ||
                `Failed to update history entry with ID ${id}`;
            console.error(`Error updating history entry ${id}:`, err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function deleteHistoryEntry(id) {
        loading.value = true;
        error.value = null;

        try {
            await historyAPI.delete(id);

            // Remove from the historyEntries array
            historyEntries.value = historyEntries.value.filter((entry) =>
                entry.id !== id
            );

            // Clear currentHistoryEntry if it's the one being deleted
            if (
                currentHistoryEntry.value &&
                currentHistoryEntry.value.id === id
            ) {
                currentHistoryEntry.value = null;
            }
        } catch (err) {
            error.value = err.message ||
                `Failed to delete history entry with ID ${id}`;
            console.error(`Error deleting history entry ${id}:`, err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    // Helper for calculating stats from history entries
    function calculateBpmProgress(exerciseId) {
        const entries = entriesByExercise.value(exerciseId);

        if (!entries.length) return [];

        // Sort by date, oldest first
        const sorted = [...entries].sort(
            (a, b) => new Date(a.start_time) - new Date(b.start_time),
        );

        // Group by day and take the highest BPM for each day
        const progressByDay = {};

        sorted.forEach((entry) => {
            const date = new Date(entry.start_time).toISOString().split("T")[0];

            if (!progressByDay[date] || entry.bpm > progressByDay[date]) {
                progressByDay[date] = entry.bpm;
            }
        });

        // Convert to array of points
        return Object.entries(progressByDay).map(([date, bpm]) => ({
            date,
            bpm,
        }));
    }

    return {
        // State
        historyEntries,
        loading,
        error,
        currentHistoryEntry,
        currentHistoryLoading,

        // Getters
        entriesSortedByDate,
        entriesByExercise,
        entriesBySession,
        entriesByDateRange,

        // Actions
        fetchHistoryEntries,
        fetchHistoryEntry,
        createHistoryEntry,
        updateHistoryEntry,
        deleteHistoryEntry,

        // Helper functions
        calculateBpmProgress,
    };
});
