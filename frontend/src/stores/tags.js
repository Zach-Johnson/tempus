import { tagsAPI } from "@/services/api.js";
import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useTagsStore = defineStore("tags", () => {
    // State
    const tags = ref([]);
    const loading = ref(false);
    const error = ref(null);
    const currentTag = ref(null);
    const currentTagLoading = ref(false);

    // Getters
    const tagsSorted = computed(
        () => {
            return [...tags.value].sort((a, b) => a.name.localeCompare(b.name));
        },
    );

    const tagById = computed(() => {
        return (id) => tags.value.find((tag) => tag.id === id);
    });

    const tagsByCategory = computed(
        () => {
            return (categoryId) =>
                tags.value.filter(
                    (tag) =>
                        tag.category_ids &&
                        tag.category_ids.includes(categoryId),
                );
        },
    );

    // Actions
    async function fetchTags(params = {}) {
        loading.value = true;
        error.value = null;

        try {
            const response = await tagsAPI.getAll(params);
            tags.value = response.data.tags || [];
        } catch (err) {
            error.value = err.message || "Failed to fetch tags";
            console.error("Error fetching tags:", err);
        } finally {
            loading.value = false;
        }
    }

    async function fetchTag(id) {
        currentTagLoading.value = true;
        error.value = null;

        try {
            const response = await tagsAPI.get(id);
            currentTag.value = response.data;

            // Also update the tag in the tags array if it exists
            const index = tags.value.findIndex((t) => t.id === id);
            if (index !== -1) {
                tags.value[index] = response.data;
            } else {
                tags.value.push(response.data);
            }
        } catch (err) {
            error.value = err.message || `Failed to fetch tag with ID ${id}`;
            console.error(`Error fetching tag ${id}:`, err);
        } finally {
            currentTagLoading.value = false;
        }
    }

    async function createTag(tagData) {
        loading.value = true;
        error.value = null;

        try {
            const response = await tagsAPI.create(tagData);
            const newTag = response.data;
            tags.value.push(newTag);
            return newTag;
        } catch (err) {
            error.value = err.message || "Failed to create tag";
            console.error("Error creating tag:", err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function updateTag(id, tagData, updateMask) {
        loading.value = true;
        error.value = null;

        try {
            // Convert paths array to comma-separated string if it's in object
            // format
            let formattedUpdateMask = updateMask;
            if (
                updateMask && typeof updateMask === "object" && updateMask.paths
            ) {
                formattedUpdateMask = updateMask.paths.join(",");
            }

            const response = await tagsAPI.update(
                id,
                tagData,
                formattedUpdateMask,
            );
            const updatedTag = response.data;

            // Update in the tags array
            const index = tags.value.findIndex((t) => t.id === id);
            if (index !== -1) {
                tags.value[index] = updatedTag;
            }

            // Update currentTag if it's the one being edited
            if (currentTag.value && currentTag.value.id === id) {
                currentTag.value = updatedTag;
            }

            return updatedTag;
        } catch (err) {
            error.value = err.message || `Failed to update tag with ID ${id}`;
            console.error(`Error updating tag ${id}:`, err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    async function deleteTag(id) {
        loading.value = true;
        error.value = null;

        try {
            await tagsAPI.delete(id);

            // Remove from the tags array
            tags.value = tags.value.filter((t) => t.id !== id);

            // Clear currentTag if it's the one being deleted
            if (currentTag.value && currentTag.value.id === id) {
                currentTag.value = null;
            }
        } catch (err) {
            error.value = err.message || `Failed to delete tag with ID ${id}`;
            console.error(`Error deleting tag ${id}:`, err);
            throw err;
        } finally {
            loading.value = false;
        }
    }

    return {
        // State
        tags,
        loading,
        error,
        currentTag,
        currentTagLoading,

        // Getters
        tagsSorted,
        tagById,
        tagsByCategory,

        // Actions
        fetchTags,
        fetchTag,
        createTag,
        updateTag,
        deleteTag,
    };
});
