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
        return practiceStats.value.total_duration_seconds /
            60; // Convert to minutes
    });

    const totalSessions = computed(() => {
        if (!practiceStats.value) return 0;
        return practiceStats.value.total_sessions;
    });

    const avgSessionTime = computed(() => {
        if (!practiceStats.value) return 0;
        return practiceStats.value.avg_session_duration_seconds /
            60; // Convert to minutes
    });

    const categoryDistribution = computed(() => {
        if (
            !practiceStats.value || !practiceStats.value.category_distribution
        ) {
            return [];
        }
        return practiceStats.value.category_distribution;
    });

    const exerciseDistribution = computed(() => {
        if (
            !practiceStats.value || !practiceStats.value.exercise_distribution
        ) {
            return [];
        }
        return practiceStats.value.exercise_distribution;
    });

    const practiceFrequency = computed(() => {
        if (!practiceStats.value || !practiceStats.value.practice_frequency) {
            return [];
        }

        // Convert practice_frequency into a format suitable for charting
        // Format for charts: { date: "YYYY-MM-DD", minutes: X, categoryId: Y }
        const chartData = [];

        practiceStats.value.practice_frequency.forEach((point) => {
            // Convert to simple date format
            const date = new Date(point.date);
            const dateStr = date.toISOString().split("T")[0]; // "YYYY-MM-DD"

            // Find category if available
            let categoryId = null;
            const minutes = point.duration_seconds / 60; // Convert to minutes

            // For now we don't have category info in practice_frequency
            // In a future version, we might add category to the backend API

            chartData.push(
                { date: dateStr, minutes: minutes, categoryId: categoryId },
            );
        });

        return chartData;
    });

    const topExercises = computed(() => {
        if (
            !practiceStats.value || !practiceStats.value.exercise_distribution
        ) {
            return [];
        }

        // Sort by duration and take top 5
        return [...practiceStats.value.exercise_distribution]
            .sort((a, b) => b.duration_seconds - a.duration_seconds)
            .slice(0, 5)
            .map((exercise) => ({
                id: exercise.exercise_id,
                name: exercise.exercise_name,
                duration: exercise.duration_seconds / 60, // Convert to minutes
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
