import { exercisesAPI, sessionsAPI } from "@/services/api";
import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useStatsStore = defineStore("stats", () => {
    // State
    const practiceStats = ref(null);
    const exerciseStats = ref({});
    const loading = ref(false);
    const error = ref(null);
    const lastRefreshed = ref(null);

    // Getters
    const totalPracticeTime = computed(() => {
        if (!practiceStats.value) return 0;
        return practiceStats.value.totalDurationSeconds /
            60; // Convert to minutes
    });

    const totalSessions = computed(() => {
        if (!practiceStats.value) return 0;
        return practiceStats.value.totalSessions;
    });

    const avgSessionTime = computed(() => {
        if (!practiceStats.value) return 0;
        return practiceStats.value.avgSessionDurationSeconds /
            60; // Convert to minutes
    });

    const categoryDistribution = computed(() => {
        if (!practiceStats.value || !practiceStats.value.categoryDistribution) {
            return [];
        }
        return practiceStats.value.categoryDistribution;
    });

    const exerciseDistribution = computed(() => {
        if (!practiceStats.value || !practiceStats.value.exerciseDistribution) {
            return [];
        }
        return practiceStats.value.exerciseDistribution;
    });

    const practiceFrequency = computed(() => {
        if (!practiceStats.value || !practiceStats.value.practiceFrequency) {
            return [];
        }

        // Convert practice_frequency into a format suitable for charting
        // Format for charts: { date: "YYYY-MM-DD", minutes: X, categoryId: Y }
        const chartData = [];

        // First add the overall practice frequency points
        practiceStats.value.practiceFrequency.forEach((point) => {
            // Convert to simple date format
            const date = new Date(point.date);
            const dateStr = date.toISOString().split("T")[0]; // "YYYY-MM-DD"
            const minutes = point.durationSeconds / 60; // Convert to minutes

            chartData.push({
                date: dateStr,
                minutes: minutes,
                categoryId: null, // null indicates overall practice time
            });
        });

        // Now add the category-specific practice frequency points if available
        if (practiceStats.value.categoryDistribution) {
            practiceStats.value.categoryDistribution.forEach((category) => {
                if (
                    category.practiceFrequency &&
                    category.practiceFrequency.length > 0
                ) {
                    category.practiceFrequency.forEach((point) => {
                        const date = new Date(point.date);
                        const dateStr = date.toISOString().split("T")[0];
                        const minutes = point.durationSeconds / 60;

                        chartData.push({
                            date: dateStr,
                            minutes: minutes,
                            categoryId: category.categoryId,
                        });
                    });
                }
            });
        }

        return chartData;
    });

    const topExercises = computed(() => {
        if (!practiceStats.value || !practiceStats.value.exerciseDistribution) {
            return [];
        }

        // Sort by duration and take top 5
        return [...practiceStats.value.exerciseDistribution]
            .sort((a, b) => b.durationSeconds - a.durationSeconds)
            .slice(0, 5)
            .map((exercise) => ({
                id: exercise.exerciseId,
                name: exercise.exerciseName,
                duration: exercise.durationSeconds / 60, // Convert to minutes
                percentage: exercise.percentage,
            }));
    });

    // Actions
    async function fetchPracticeStats(params = {}) {
        loading.value = true;
        error.value = null;

        try {
            const response = await sessionsAPI.getStats(params);
            practiceStats.value = response.data;
            lastRefreshed.value = new Date();
            return response.data;
        } catch (err) {
            error.value = err.message || "Failed to fetch practice statistics";
            console.error("Error fetching practice statistics:", err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function fetchExerciseStats(exerciseId, params = {}) {
        if (!exerciseId) return null;

        loading.value = true;
        error.value = null;

        try {
            const response = await exercisesAPI.getStats(exerciseId, params);
            exerciseStats.value[exerciseId] = response.data;
            return response.data;
        } catch (err) {
            error.value = err.message ||
                `Failed to fetch stats for exercise ID ${exerciseId}`;
            console.error(
                `Error fetching exercise stats for ID ${exerciseId}:`,
                err,
            );
            throw err;
        } finally {
            loading.value = false;
        }
    }

    // Return the store
    return {
        // State
        practiceStats,
        exerciseStats,
        loading,
        error,
        lastRefreshed,

        // Getters
        totalPracticeTime,
        totalSessions,
        avgSessionTime,
        categoryDistribution,
        exerciseDistribution,
        practiceFrequency,
        topExercises,

        // Actions
        fetchPracticeStats,
        fetchExerciseStats,
    };
});
