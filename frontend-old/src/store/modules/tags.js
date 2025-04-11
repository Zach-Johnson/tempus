import tagService from "@/services/tagService";
import { defineStore } from "pinia";

export const useTagsStore = defineStore("tags", {
    state: () => ({ tags: [], totalCount: 0, loading: false, error: null }),

    getters: {
        getTagById: (state) => (id) => {
            return state.tags.find((tag) => tag.id === id);
        },

        tagsOptions: (state) => {
            return state.tags.map((tag) => ({
                value: tag.id,
                label: tag.name,
            }));
        },
    },

    actions: {
        async fetchTags(params = {
            page: 1,
            page_size: 10,
        }) {
            this.loading = true;
            this.error = null;

            try {
                const response = await tagService.getTags(params);
                this.tags = response.tags || [];
                this.totalCount = response.total_count || 0;
            } catch (error) {
                this.error = error.message || "Failed to fetch tags";
                console.error("Error fetching tags:", error);
            } finally {
                this.loading = false;
            }
        },

        async fetchTagById(id) {
            this.loading = true;
            this.error = null;

            try {
                const tag = await tagService.getTag(id);
                return tag;
            } catch (error) {
                this.error = error.message || `Failed to fetch tag ${id}`;
                console.error(`Error fetching tag ${id}:`, error);
                return null;
            } finally {
                this.loading = false;
            }
        },

        async createTag(tagData) {
            this.loading = true;
            this.error = null;

            try {
                const tag = await tagService.createTag(tagData);
                this.tags.push(tag);
                this.totalCount++;
                return tag;
            } catch (error) {
                this.error = error.message || "Failed to create tag";
                console.error("Error creating tag:", error);
                return null;
            } finally {
                this.loading = false;
            }
        },

        async deleteTag(id) {
            this.loading = true;
            this.error = null;

            try {
                await tagService.deleteTag(id);

                // Remove from the list
                this.tags = this.tags.filter((t) => t.id !== id);
                this.totalCount--;

                return true;
            } catch (error) {
                this.error = error.message || `Failed to delete tag ${id}`;
                console.error(`Error deleting tag ${id}:`, error);
                return false;
            } finally {
                this.loading = false;
            }
        },
    },
});
