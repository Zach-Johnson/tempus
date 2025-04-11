import { exercisesAPI } from "@/services/api.js";
import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useExercisesStore = defineStore("exercises", () => {
    // State
    const exercises = ref([]);
    const loading = ref(false);
    const error = ref(null);
    const currentExercise = ref(null);
    const currentExerciseLoading = ref(false);
    const exerciseStats = ref(null);
    const statsLoading = ref(false);

    // Getters
    const exercisesSorted = computed(
        () => {
            return [...exercises.value].sort(
                (a, b) => a.name.localeCompare(b.name),
            );
        },
    );

    const exerciseById = computed(() => {
        return (id) => exercises.value.find((ex) => ex.id === id);
    });

    const exercisesByCategory = computed(
        () => {
            return (categoryId) =>
                exercises.value.filter(
                    (ex) =>
                        ex.category_ids &&
                        ex.category_ids.includes(categoryId),
                );
        },
    );

    const exercisesByTag = computed(
        () => {
            return (tagId) =>
                exercises.value.filter(
                    (ex) => ex.tag_ids && ex.tag_ids.includes(tagId),
                );
        },
    );

    // Actions
    async function fetchExercises(params = {}) {
        loading.value = true;
        error.value = null;

        try {
            const response = await exercisesAPI.getAll(params);
            exercises.value = response.data.exercises || [];
        } catch (err) {
            error.value = err.message || "Failed to fetch exercises";
            console.error("Error fetching exercises:", err);
        } finally {
            loading.value = false;
        }
    }

    async function fetchExercise(id) {
        currentExerciseLoading.value = true;
        error.value = null;

        try {
            const response = await exercisesAPI.get(id);
            currentExercise.value = response.data;

            // Also update the exercise in the exercises array if it exists
            const index = exercises.value.findIndex((e) => e.id === id);
            if (index !== -1) {
                exercises.value[index] = response.data;
            } else {
                exercises.value.push(response.data);
            }
        } catch (err) {
            error.value = err.message ||
                `Failed to fetch exercise with ID ${id}`;
            console.error(`Error fetching exercise ${id}:`, err);
        } finally {
            currentExerciseLoading.value = false;
        }
    }

    async function createExercise(exerciseData) {
        loading.value = true;
        error.value = null;

        try {
            const response = await exercisesAPI.create(exerciseData);
            const newExercise = response.data;
            exercises.value.push(newExercise);
            return newExercise;
        } catch (err) {
            error.value = err.message || "Failed to create exercise";
            console.error("Error creating exercise:", err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function updateExercise(id, exerciseData, updateMask) {
        loading.value = true;
        error.value = null;

        try {
            const response = await exercisesAPI.update(
                id,
                exerciseData,
                updateMask,
            );
            const updatedExercise = response.data;

            // Update in the exercises array
            const index = exercises.value.findIndex((e) => e.id === id);
            if (index !== -1) {
                exercises.value[index] = updatedExercise;
            }

            // Update currentExercise if it's the one being edited
            if (currentExercise.value && currentExercise.value.id === id) {
                currentExercise.value = updatedExercise;
            }

            return updatedExercise;
        } catch (err) {
            error.value = err.message ||
                `Failed to update exercise with ID ${id}`;
            console.error(`Error updating exercise ${id}:`, err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function deleteExercise(id) {
        loading.value = true;
        error.value = null;

        try {
            await exercisesAPI.delete(id);

            // Remove from the exercises array
            exercises.value = exercises.value.filter((e) => e.id !== id);

            // Clear currentExercise if it's the one being deleted
            if (currentExercise.value && currentExercise.value.id === id) {
                currentExercise.value = null;
            }
        } catch (err) {
            error.value = err.message ||
                `Failed to delete exercise with ID ${id}`;
            console.error(`Error deleting exercise ${id}:`, err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function addExerciseImage(exerciseId, imageData) {
        loading.value = true;
        error.value = null;

        try {
            const response = await exercisesAPI.addImage(exerciseId, imageData);

            // Update exercise with new image if it's in the exercises array
            const exercise = exercises.value.find((e) => e.id === exerciseId);
            if (exercise) {
                if (!exercise.images) {
                    exercise.images = [];
                }
                exercise.images.push(response.data);
            }

            // Update currentExercise if it's the relevant exercise
            if (
                currentExercise.value &&
                currentExercise.value.id === exerciseId
            ) {
                if (!currentExercise.value.images) {
                    currentExercise.value.images = [];
                }
                currentExercise.value.images.push(response.data);
            }

            return response.data;
        } catch (err) {
            error.value = err.message ||
                `Failed to add image to exercise with ID ${exerciseId}`;
            console.error(`Error adding image to exercise ${exerciseId}:`, err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function deleteExerciseImage(imageId) {
        loading.value = true;
        error.value = null;

        try {
            await exercisesAPI.deleteImage(imageId);

            // Remove image from exercises in the array
            exercises.value.forEach((exercise) => {
                if (exercise.images) {
                    exercise.images = exercise.images.filter((img) =>
                        img.id !== imageId
                    );
                }
            });

            // Remove image from currentExercise if applicable
            if (currentExercise.value && currentExercise.value.images) {
                currentExercise.value.images = currentExercise.value.images
                    .filter(
                        (img) => img.id !== imageId,
                    );
            }
        } catch (err) {
            error.value = err.message ||
                `Failed to delete image with ID ${imageId}`;
            console.error(`Error deleting image ${imageId}:`, err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function addExerciseLink(exerciseId, linkData) {
        loading.value = true;
        error.value = null;

        try {
            const response = await exercisesAPI.addLink(exerciseId, linkData);

            // Update exercise with new link if it's in the exercises array
            const exercise = exercises.value.find((e) => e.id === exerciseId);
            if (exercise) {
                if (!exercise.links) {
                    exercise.links = [];
                }
                exercise.links.push(response.data);
            }

            // Update currentExercise if it's the relevant exercise
            if (
                currentExercise.value &&
                currentExercise.value.id === exerciseId
            ) {
                if (!currentExercise.value.links) {
                    currentExercise.value.links = [];
                }
                currentExercise.value.links.push(response.data);
            }

            return response.data;
        } catch (err) {
            error.value = err.message ||
                `Failed to add link to exercise with ID ${exerciseId}`;
            console.error(`Error adding link to exercise ${exerciseId}:`, err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function deleteExerciseLink(linkId) {
        loading.value = true;
        error.value = null;

        try {
            await exercisesAPI.deleteLink(linkId);

            // Remove link from exercises in the array
            exercises.value.forEach((exercise) => {
                if (exercise.links) {
                    exercise.links = exercise.links.filter((link) =>
                        link.id !== linkId
                    );
                }
            });

            // Remove link from currentExercise if applicable
            if (currentExercise.value && currentExercise.value.links) {
                currentExercise.value.links = currentExercise.value.links
                    .filter(
                        (link) => link.id !== linkId,
                    );
            }
        } catch (err) {
            error.value = err.message ||
                `Failed to delete link with ID ${linkId}`;
            console.error(`Error deleting link ${linkId}:`, err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function fetchExerciseStats(exerciseId, params = {}) {
        statsLoading.value = true;
        error.value = null;

        try {
            const response = await exercisesAPI.getStats(exerciseId, params);
            exerciseStats.value = response.data;
            return response.data;
        } catch (err) {
            error.value = err.message ||
                `Failed to fetch stats for exercise with ID ${exerciseId}`;
            console.error(
                `Error fetching stats for exercise ${exerciseId}:`,
                err,
            );
            throw err;
        } finally {
            statsLoading.value = false;
        }
    }

    return {
        // State
        exercises,
        loading,
        error,
        currentExercise,
        currentExerciseLoading,
        exerciseStats,
        statsLoading,

        // Getters
        exercisesSorted,
        exerciseById,
        exercisesByCategory,
        exercisesByTag,

        // Actions
        fetchExercises,
        fetchExercise,
        createExercise,
        updateExercise,
        deleteExercise,
        addExerciseImage,
        deleteExerciseImage,
        addExerciseLink,
        deleteExerciseLink,
        fetchExerciseStats,
    };
});
