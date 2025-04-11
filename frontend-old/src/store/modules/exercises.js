import exerciseService from "@/services/exerciseService";
import { defineStore } from "pinia";

export const useExercisesStore = defineStore("exercises", {
    state: () => ({
        exercises: [],
        currentExercise: null,
        totalCount: 0,
        loading: false,
        error: null,
        filters: { categoryIds: [], tagIds: [], searchTerm: "" },
    }),

    getters: {
        getExerciseById: (state) => (id) => {
            return state.exercises.find((exercise) => exercise.id === id);
        },

        exercisesOptions: (state) => {
            return state.exercises.map(
                (exercise) => ({ value: exercise.id, label: exercise.name }),
            );
        },
    },

    actions: {
        setFilters(filters) {
            this.filters = { ...this.filters, ...filters };
        },

        clearFilters() {
            this.filters = { categoryIds: [], tagIds: [], searchTerm: "" };
        },

        async fetchExercises(params = {
            page: 1,
            page_size: 10,
        }) {
            this.loading = true;
            this.error = null;

            try {
                // Apply filters
                const queryParams = {
                    ...params,
                    category_ids: this.filters.categoryIds,
                    tag_ids: this.filters.tagIds,
                    search_term: this.filters.searchTerm,
                };

                const response = await exerciseService.getExercises(
                    queryParams,
                );
                this.exercises = response.exercises || [];
                this.totalCount = response.total_count || 0;
            } catch (error) {
                this.error = error.message || "Failed to fetch exercises";
                console.error("Error fetching exercises:", error);
            } finally {
                this.loading = false;
            }
        },

        async fetchExerciseById(id) {
            this.loading = true;
            this.error = null;

            try {
                const exercise = await exerciseService.getExercise(id);
                this.currentExercise = exercise;
                return exercise;
            } catch (error) {
                this.error = error.message || `Failed to fetch exercise ${id}`;
                console.error(`Error fetching exercise ${id}:`, error);
                return null;
            } finally {
                this.loading = false;
            }
        },

        async createExercise(exerciseData) {
            this.loading = true;
            this.error = null;

            try {
                const exercise = await exerciseService.createExercise(
                    exerciseData,
                );
                this.exercises.push(exercise);
                this.totalCount++;
                return exercise;
            } catch (error) {
                this.error = error.message || "Failed to create exercise";
                console.error("Error creating exercise:", error);
                return null;
            } finally {
                this.loading = false;
            }
        },

        async updateExercise(exerciseData) {
            this.loading = true;
            this.error = null;

            try {
                const updatedExercise = await exerciseService.updateExercise(
                    exerciseData,
                );

                // Update in the list
                const index = this.exercises.findIndex((e) =>
                    e.id === updatedExercise.id
                );
                if (index !== -1) {
                    this.exercises[index] = updatedExercise;
                }

                // Update current exercise if it's the same one
                if (
                    this.currentExercise &&
                    this.currentExercise.id === updatedExercise.id
                ) {
                    this.currentExercise = updatedExercise;
                }

                return updatedExercise;
            } catch (error) {
                this.error = error.message || "Failed to update exercise";
                console.error("Error updating exercise:", error);
                return null;
            } finally {
                this.loading = false;
            }
        },

        async deleteExercise(id) {
            this.loading = true;
            this.error = null;

            try {
                await exerciseService.deleteExercise(id);

                // Remove from the list
                this.exercises = this.exercises.filter((e) => e.id !== id);
                this.totalCount--;

                // Clear current exercise if it's the same one
                if (this.currentExercise && this.currentExercise.id === id) {
                    this.currentExercise = null;
                }

                return true;
            } catch (error) {
                this.error = error.message || `Failed to delete exercise ${id}`;
                console.error(`Error deleting exercise ${id}:`, error);
                return false;
            } finally {
                this.loading = false;
            }
        },

        async addExerciseLink(link) {
            this.loading = true;
            this.error = null;

            try {
                const newLink = await exerciseService.addExerciseLink(link);

                // Update the exercise links if it's the current exercise
                if (
                    this.currentExercise &&
                    this.currentExercise.id === link.exercise_id
                ) {
                    if (!this.currentExercise.links) {
                        this.currentExercise.links = [];
                    }
                    this.currentExercise.links.push(newLink);
                }

                return newLink;
            } catch (error) {
                this.error = error.message || "Failed to add exercise link";
                console.error("Error adding exercise link:", error);
                return null;
            } finally {
                this.loading = false;
            }
        },

        async deleteExerciseLink(id) {
            this.loading = true;
            this.error = null;

            try {
                await exerciseService.deleteExerciseLink(id);

                // Remove from the current exercise if applicable
                if (this.currentExercise && this.currentExercise.links) {
                    this.currentExercise.links = this.currentExercise.links
                        .filter(
                            (link) => link.id !== id,
                        );
                }

                return true;
            } catch (error) {
                this.error = error.message ||
                    `Failed to delete exercise link ${id}`;
                console.error(`Error deleting exercise link ${id}:`, error);
                return false;
            } finally {
                this.loading = false;
            }
        },

        async addExerciseToCategory(exerciseId, categoryId) {
            this.loading = true;
            this.error = null;

            try {
                await exerciseService.addExerciseToCategory(
                    exerciseId,
                    categoryId,
                );

                // Refetch the exercise to get updated categories
                if (
                    this.currentExercise &&
                    this.currentExercise.id === exerciseId
                ) {
                    await this.fetchExerciseById(exerciseId);
                }

                return true;
            } catch (error) {
                this.error = error.message ||
                    "Failed to add exercise to category";
                console.error("Error adding exercise to category:", error);
                return false;
            } finally {
                this.loading = false;
            }
        },

        async removeExerciseFromCategory(exerciseId, categoryId) {
            this.loading = true;
            this.error = null;

            try {
                await exerciseService.removeExerciseFromCategory(
                    exerciseId,
                    categoryId,
                );

                // Update the current exercise if applicable
                if (
                    this.currentExercise &&
                    this.currentExercise.id === exerciseId
                ) {
                    this.currentExercise.categories = this.currentExercise
                        .categories.filter(
                            (category) => category.id !== categoryId,
                        );
                }

                return true;
            } catch (error) {
                this.error = error.message ||
                    "Failed to remove exercise from category";
                console.error("Error removing exercise from category:", error);
                return false;
            } finally {
                this.loading = false;
            }
        },

        async addTagToExercise(exerciseId, tagId) {
            this.loading = true;
            this.error = null;

            try {
                await exerciseService.addTagToExercise(exerciseId, tagId);

                // Refetch the exercise to get updated tags
                if (
                    this.currentExercise &&
                    this.currentExercise.id === exerciseId
                ) {
                    await this.fetchExerciseById(exerciseId);
                }

                return true;
            } catch (error) {
                this.error = error.message || "Failed to add tag to exercise";
                console.error("Error adding tag to exercise:", error);
                return false;
            } finally {
                this.loading = false;
            }
        },

        async removeTagFromExercise(exerciseId, tagId) {
            this.loading = true;
            this.error = null;

            try {
                await exerciseService.removeTagFromExercise(exerciseId, tagId);

                // Update the current exercise if applicable
                if (
                    this.currentExercise &&
                    this.currentExercise.id === exerciseId
                ) {
                    this.currentExercise.tags = this.currentExercise.tags
                        .filter(
                            (tag) => tag.id !== tagId,
                        );
                }

                return true;
            } catch (error) {
                this.error = error.message ||
                    "Failed to remove tag from exercise";
                console.error("Error removing tag from exercise:", error);
                return false;
            } finally {
                this.loading = false;
            }
        },
    },
});
