import historyService from "@/services/historyService";
import { defineStore } from "pinia";

export const useHistoryStore = defineStore("history", {
    state: () => ({
        exerciseHistory: [],
        totalCount: 0,
        loading: false,
        error: null,
        currentExerciseId: null,
        dateRange: { startDate: null, endDate: null },
    }),

    getters: {
        getHistoryById: (state) => (id) => {
            return state.exerciseHistory.find((entry) => entry.id === id);
        },

        historyByDate: (state) => {
            const historyMap = {};

            state.exerciseHistory.forEach((entry) => {
                const date = new Date(entry.date).toISOString().split("T")[0];
                if (!historyMap[date]) {
                    historyMap[date] = [];
                }
                historyMap[date].push(entry);
            });

            return historyMap;
        },

        averageBpm: (state) => {
            if (state.exerciseHistory.length === 0) return 0;

            const totalBpm = state.exerciseHistory.reduce((sum, entry) => {
                return sum + (entry.bpm || 0);
            }, 0);

            return Math.round(totalBpm / state.exerciseHistory.length);
        },

        averageRating: (state) => {
            if (state.exerciseHistory.length === 0) return 0;

            const totalRating = state.exerciseHistory.reduce((sum, entry) => {
                return sum + (entry.rating || 0);
            }, 0);

            return (totalRating / state.exerciseHistory.length).toFixed(1);
        },

        progressData: (state) => {
            // Group history by date and calculate average BPMs
            const groupedData = {};

            state.exerciseHistory.forEach((entry) => {
                if (!entry.bpm) return;

                const date = new Date(entry.date).toISOString().split("T")[0];
                if (!groupedData[date]) {
                    groupedData[date] = { sum: 0, count: 0 };
                }

                groupedData[date].sum += entry.bpm;
                groupedData[date].count += 1;
            });

            // Convert to array sorted by date
            return Object.keys(groupedData)
                .sort()
                .map((date) => ({
                    date,
                    bpm: Math.round(
                        groupedData[date].sum / groupedData[date].count,
                    ),
                }));
        },
    },

    actions: {
        setExerciseId(exerciseId) {
            this.currentExerciseId = exerciseId;
        },

        setDateRange(startDate, endDate) {
            this.dateRange = { startDate, endDate };
        },

        clearDateRange() {
            this.dateRange = { startDate: null, endDate: null };
        },

        async fetchExerciseHistory(params = {
            page: 1,
            page_size: 10,
        }) {
            this.loading = true;
            this.error = null;

            try {
                // Apply filters
                const queryParams = {
                    ...params,
                    exercise_id: this.currentExerciseId,
                    start_date: this.dateRange.startDate,
                    end_date: this.dateRange.endDate,
                };

                const response = await historyService.getExerciseHistory(
                    queryParams,
                );
                this.exerciseHistory = response.history_entries || [];
                this.totalCount = response.total_count || 0;
            } catch (error) {
                this.error = error.message ||
                    "Failed to fetch exercise history";
                console.error("Error fetching exercise history:", error);
            } finally {
                this.loading = false;
            }
        },

        async createExerciseHistory(historyData) {
            this.loading = true;
            this.error = null;

            try {
                const historyEntry = await historyService.createExerciseHistory(
                    historyData,
                );
                this.exerciseHistory.push(historyEntry);
                this.totalCount++;
                return historyEntry;
            } catch (error) {
                this.error = error.message ||
                    "Failed to create exercise history";
                console.error("Error creating exercise history:", error);
                return null;
            } finally {
                this.loading = false;
            }
        },

        async deleteExerciseHistory(id) {
            this.loading = true;
            this.error = null;

            try {
                await historyService.deleteExerciseHistory(id);

                // Remove from the list
                this.exerciseHistory = this.exerciseHistory.filter((h) =>
                    h.id !== id
                );
                this.totalCount--;

                return true;
            } catch (error) {
                this.error = error.message ||
                    `Failed to delete exercise history ${id}`;
                console.error(`Error deleting exercise history ${id}:`, error);
                return false;
            } finally {
                this.loading = false;
            }
        },
    },
});
